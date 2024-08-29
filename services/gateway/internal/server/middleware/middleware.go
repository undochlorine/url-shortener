package middleware

import "net/http"

type Middleware func(handler http.Handler) http.HandlerFunc

func Chain(chain ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for _, c := range chain {
			next = c(next)
		}

		return next.ServeHTTP
	}
}
