package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/danielgtaylor/huma/v2"

	"github.com/rlibaert/service-example-go/domain"
	"github.com/rlibaert/service-example-go/restapi"
	"github.com/rlibaert/service-example-go/router"
	"github.com/rlibaert/service-example-go/stores"
	"github.com/rlibaert/service-example-go/wrappers"
)

type ServerOptions struct {
	Host              string        `short:"H" doc:"host to listen on"                    default:""`
	Port              string        `short:"p" doc:"port to listen on"                    default:"8888"`
	ReadHeaderTimeout time.Duration `          doc:"time allowed to read request headers" default:"15s"`
}

func NewServer(options *ServerOptions, handler http.Handler, logger *slog.Logger) *http.Server {
	return &http.Server{
		Addr:              options.Host + ":" + options.Port,
		ReadHeaderTimeout: options.ReadHeaderTimeout,
		Handler:           handler,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
}

type RouterOptions struct {
	EndpointsPrefix string `doc:"mount endpoints at a prefix" default:"/api"`
}

func NewRouter(
	options *RouterOptions,
	title string,
	version string,
	revision string,
	created string,
	logger *slog.Logger,
) http.Handler {
	buildinfoMetric := joinQuote("build_info{goversion=", runtime.Version(),
		",title=", title,
		",version=", version,
		",revision=", revision,
		",created=", created,
		"} 1\n")
	metriks := metrics.NewSet()
	return router.New(title, version,
		func(_ http.ResponseWriter, _ *http.Request) {},
		func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprint(w, buildinfoMetric)
			metriks.WritePrometheus(w)
			metrics.WriteProcessMetrics(w)
		},
		router.OptUseMiddleware(
			ctxlog{}.loggerMiddleware(logger),
			router.RequestsMetricsMiddleware(metriks),
			router.ResponsesMetricsMiddleware(metriks),
			ctxlog{}.recoverMiddleware(logger),
		),
		router.OptGroup(options.EndpointsPrefix,
			router.OptAutoRegister(&restapi.ServiceRegisterer{
				Service: wrappers.ServiceErrorHandler{
					Service: &domain.ServiceStore{
						Store: stores.MustNewMock(&domain.Contact{
							Firstname: "john",
							Lastname:  "smith",
							Birthday:  time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC),
						}),
					},
					ErrorHandler: func(ctx context.Context, err error) {
						ctxlog{}.get(ctx).
							LogAttrs(context.Background(), slog.LevelError, "service error", slog.Any("err", err))
					},
				},
			}),
			router.OptAutoRegister(&restapi.GreetRegisterer{}),
			router.OptAutoRegister(&restapi.PanicRegisterer{}),
		),
	)
}

// ctxlog is a [context.Context] key and acts as a virtual package for operations related to it.
type ctxlog struct{}

// loggerMiddleware returns a middleware that sets a [slog.Logger] in
// the [context.Context] and logs the request after it has terminated.
func (key ctxlog) loggerMiddleware(parent *slog.Logger) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		logger := parent.With("x-request-id", ctx.Header("X-Request-Id"))

		start := time.Now()
		next(huma.WithValue(ctx, key, logger.WithGroup("op").With("id", ctx.Operation().OperationID)))

		logger.LogAttrs(context.Background(), slog.LevelInfo,
			joinSpace(ctx.Operation().Method, ctx.Operation().Path, ctx.Version().Proto),
			slog.String("from", ctx.RemoteAddr()),
			slog.String("ref", ctx.Header("Referer")),
			slog.String("ua", ctx.Header("User-Agent")),
			slog.Int("status", ctx.Status()),
			slog.Duration("dur", time.Since(start)),
		)
	}
}

func (key ctxlog) get(ctx context.Context) *slog.Logger {
	l, _ := ctx.Value(key).(*slog.Logger)
	return l
}

// recoverMiddleware returns a middleware that recovers and logs the value from panic
// to finally set the response status to [http.StatusInternalServerError].
func (key ctxlog) recoverMiddleware(fallback *slog.Logger) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		defer func() {
			v := recover()
			if v != nil {
				logger, ok := ctx.Context().Value(key).(*slog.Logger)
				if !ok {
					logger = fallback
				}
				logger.LogAttrs(context.Background(), slog.LevelError, "panic occurred", slog.Any("recovered", v))
				ctx.SetStatus(http.StatusInternalServerError)
			}
		}()
		next(ctx)
	}
}

// joinQuote is [strings.Join] with " as separator.
func joinQuote(elems ...string) string { return strings.Join(elems, `"`) }

// joinSpace is [strings.Join] with space as separator.
func joinSpace(elems ...string) string { return strings.Join(elems, ` `) }
