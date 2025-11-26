// Package router provides domain-agnostic primitives to build [huma]-based routers.
package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

// New returns a new [huma]-based router.
func New(
	title, version string,
	readiness http.HandlerFunc,
	metrics http.HandlerFunc,
	opts ...func(huma.API),
) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/liveness", func(http.ResponseWriter, *http.Request) {})
	mux.HandleFunc("/readiness", readiness)
	mux.HandleFunc("/metrics", metrics)

	api := humago.New(mux, huma.DefaultConfig(title, version))
	for _, opt := range opts {
		opt(api)
	}

	return mux
}

// OptUseMiddleware returns a [huma.API] option to append new middlewares.
func OptUseMiddleware(middlewares ...func(huma.Context, func(huma.Context))) func(huma.API) {
	return func(api huma.API) { api.UseMiddleware(middlewares...) }
}

// OptAutoRegister returns a [huma.API] option to auto-detect and call a server registration methods.
func OptAutoRegister(server any) func(huma.API) {
	return func(api huma.API) { huma.AutoRegister(api, server) }
}

type Prefixes []string

// OptGroup returns a [huma.API] option that creates a new group to apply options.
func (p Prefixes) OptGroup(opts ...func(huma.API)) func(huma.API) {
	return func(api huma.API) {
		g := huma.NewGroup(api, p...)
		for _, opt := range opts {
			opt(g)
		}
	}
}

// OptGroup returns a [huma.API] option that creates a new group at a prefix to apply options.
func OptGroup(prefix string, opts ...func(huma.API)) func(huma.API) {
	return Prefixes{prefix}.OptGroup(opts...)
}
