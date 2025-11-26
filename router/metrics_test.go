package router_test

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"

	"github.com/rlibaert/service-example-go/router"
)

func ExampleRequestsMetricsMiddleware() {
	set := metrics.NewSet()
	handler := huma.Middlewares{router.RequestsMetricsMiddleware(set)}.
		Handler(func(huma.Context) {})
	op := huma.Operation{Method: http.MethodGet, Path: "/teapot"}

	handler(humatest.NewContext(&op, nil, nil))
	set.WritePrometheus(os.Stdout)

	// Output:
	// http_requests_in_flight{method="GET",path="/teapot"} 0
}

func ExampleRequestsMetricsMiddleware_inflight() {
	set := metrics.NewSet()
	handler := huma.Middlewares{router.RequestsMetricsMiddleware(set)}.
		Handler(func(huma.Context) { set.WritePrometheus(os.Stdout) })
	op := huma.Operation{Method: http.MethodGet, Path: "/teapot"}

	handler(humatest.NewContext(&op, nil, nil))

	// Output:
	// http_requests_in_flight{method="GET",path="/teapot"} 1
}

func ExampleResponsesMetricsMiddleware() {
	set := metrics.NewSet()
	handler := huma.Middlewares{router.ResponsesMetricsMiddleware(set)}.
		Handler(func(ctx huma.Context) { ctx.SetStatus(http.StatusTeapot) })
	op := huma.Operation{Method: http.MethodGet, Path: "/teapot"}

	for range 42 {
		handler(humatest.NewContext(&op, nil, httptest.NewRecorder()))
	}

	var buf bytes.Buffer
	set.WritePrometheus(&buf)
	scanner := bufio.NewScanner(&buf)

	for range 7 {
		scanner.Scan()
		fmt.Println(scanner.Text())
	}

	scanner.Scan()
	sum, _, ok := strings.Cut(scanner.Text(), " ")
	fmt.Println(sum, "value skipped", ok)

	// Output:
	// http_request_duration_seconds_bucket{method="GET",path="/teapot",status="418",le="0.001"} 42
	// http_request_duration_seconds_bucket{method="GET",path="/teapot",status="418",le="0.005"} 42
	// http_request_duration_seconds_bucket{method="GET",path="/teapot",status="418",le="0.025"} 42
	// http_request_duration_seconds_bucket{method="GET",path="/teapot",status="418",le="0.125"} 42
	// http_request_duration_seconds_bucket{method="GET",path="/teapot",status="418",le="0.625"} 42
	// http_request_duration_seconds_bucket{method="GET",path="/teapot",status="418",le="3.125"} 42
	// http_request_duration_seconds_bucket{method="GET",path="/teapot",status="418",le="+Inf"} 42
	// http_request_duration_seconds_sum{method="GET",path="/teapot",status="418"} value skipped true
}

func BenchmarkMetrics(b *testing.B) {
	set := metrics.NewSet()
	handler := huma.Middlewares{
		router.RequestsMetricsMiddleware(set),
		router.ResponsesMetricsMiddleware(set),
	}.Handler(func(huma.Context) {})
	op := huma.Operation{Method: http.MethodGet, Path: "/teapot"}

	for b.Loop() {
		handler(humatest.NewContext(&op, nil, nil))
	}

	if d := b.Elapsed() / time.Duration(b.N); d > time.Microsecond {
		b.Error(b.Name(), "is too slow: took", d, "per op")
	}
}
