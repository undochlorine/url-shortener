package middleware

import "net/http"

func CORS(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, Accept, Origin")
			w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")

			return
		}

		next.ServeHTTP(w, r)
	}
}
