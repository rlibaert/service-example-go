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
			ctxlog{}.setMiddleware(logger),
			router.RequestsLogMiddleware(func(ctx context.Context, r slog.Record) {
				ctxlog{}.get(ctx).Handler().Handle(ctx, r) //nolint: errcheck,gosec // ignored by [slog.Logger.Log] as well
			}),
			router.RequestsMetricsMiddleware(metriks),
			router.ResponsesMetricsMiddleware(metriks),
			router.RecoverMiddleware(func(ctx context.Context, a any) {
				ctxlog{}.get(ctx).LogAttrs(ctx, slog.LevelError, "panic occurred", slog.Any("recovered", a))
			}),
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

func (key ctxlog) setMiddleware(parent *slog.Logger) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		logger := parent.With("x-request-id", ctx.Header("X-Request-Id"))
		ctx = huma.WithValue(ctx, key, logger)
		next(ctx)
	}
}

func (key ctxlog) get(ctx context.Context) *slog.Logger {
	l, _ := ctx.Value(key).(*slog.Logger)
	return l
}

// joinQuote is [strings.Join] with " as separator.
func joinQuote(elems ...string) string { return strings.Join(elems, `"`) }
