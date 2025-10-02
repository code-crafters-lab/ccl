package main

import (
	"net/http"

	"github.com/rs/cors"
)

var (
	DefaultCORSOptions = cors.Options{
		AllowCredentials: true,
		AllowedHeaders: []string{
			"*",
		},
		AllowedMethods: []string{
			http.MethodOptions,
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		ExposedHeaders: []string{
			"location",
			"content-length",
			"grpc-status",
			"grpc-message",
			"grpc-status-details-bin",
		},
		AllowOriginFunc: func(_ string) bool {
			return true
		},
	}
)

func CORSInterceptorOpts(opts cors.Options, h http.Handler) http.Handler {
	return cors.New(opts).Handler(h)
}

func CORSInterceptor(h http.Handler) http.Handler {
	return CORSInterceptorOpts(DefaultCORSOptions, h)
}
