package router

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/danielgtaylor/huma/v2"
)

// RequestsLogMiddleware creates a [slog.Record] for done requests and calls a handling function.
func RequestsLogMiddleware(handle func(context.Context, slog.Record)) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		start := time.Now()
		defer func() {
			msg := joinSpace(ctx.Method(), ctx.URL().Path, ctx.Version().Proto)
			rec := slog.NewRecord(time.Now(), slog.LevelInfo, msg, 0)
			rec.AddAttrs(
				slog.String("from", ctx.RemoteAddr()),
				slog.String("ref", ctx.Header("Referer")),
				slog.String("ua", ctx.Header("User-Agent")),
				slog.Int("status", ctx.Status()),
				slog.Duration("dur", rec.Time.Sub(start)),
			)
			handle(ctx.Context(), rec)
		}()
		next(ctx)
	}
}

// RequestsMetricsMiddleware returns a middleware collecting requests metrics.
//
//   - http_requests_in_flight{method,path}
func RequestsMetricsMiddleware(set *metrics.Set) func(huma.Context, func(huma.Context)) {
	var m sync.Map
	return func(ctx huma.Context, next func(huma.Context)) {
		op := ctx.Operation()
		k := op.OperationID
		v, ok := m.Load(k)
		if !ok {
			labels := joinQuote("{method=", op.Method, ",path=", op.Path, "}")
			v, _ = m.LoadOrStore(k,
				set.GetOrCreateCounter("http_requests_in_flight"+labels),
			)
		}
		val := v.(*metrics.Counter) //nolint: errcheck // always true
		val.Inc()
		defer val.Dec()

		next(ctx)
	}
}

// ResponsesMetricsMiddleware returns a middleware collecting request responses metrics.
//
//   - http_request_duration_seconds_bucket{method,path,status,le}
//   - http_request_duration_seconds_sum{method,path,status}
//   - http_request_duration_seconds_count{method,path,status}
//   - http_requests_total{method,path,status}
func ResponsesMetricsMiddleware(set *metrics.Set) func(huma.Context, func(huma.Context)) {
	type value struct {
		*metrics.PrometheusHistogram
		*metrics.Counter
	}
	var buckets = metrics.ExponentialBuckets(1e-3, 5, 6) //nolint: mnd // arbitrary

	var m sync.Map
	return func(ctx huma.Context, next func(huma.Context)) {
		start := time.Now()
		defer func() {
			op := ctx.Operation()
			k := op.OperationID + http.StatusText(ctx.Status())
			v, ok := m.Load(k)
			if !ok {
				labels := joinQuote("{method=", op.Method, ",path=", op.Path, ",status=", strconv.Itoa(ctx.Status()), "}") //nolint: golines
				v, _ = m.LoadOrStore(k, value{
					set.GetOrCreatePrometheusHistogramExt("http_request_duration_seconds"+labels, buckets),
					set.GetOrCreateCounter("http_requests_total" + labels),
				})
			}
			val := v.(value) //nolint: errcheck // always true
			val.PrometheusHistogram.UpdateDuration(start)
			val.Counter.Inc()
		}()

		next(ctx)
	}
}

// RecoverMiddleware returns a middleware to [recover] from [panic] and sets the
// response status to [http.StatusInternalServerError] and handles the recovered value.
func RecoverMiddleware(handle func(context.Context, any)) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		defer func() {
			a := recover()
			if a != nil {
				ctx.SetStatus(http.StatusInternalServerError)
				handle(ctx.Context(), a)
			}
		}()
		next(ctx)
	}
}

// joinQuote is [strings.Join] with " as separator.
func joinQuote(elems ...string) string { return strings.Join(elems, `"`) }

// joinSpace is [strings.Join] with space as separator.
func joinSpace(elems ...string) string { return strings.Join(elems, ` `) }
