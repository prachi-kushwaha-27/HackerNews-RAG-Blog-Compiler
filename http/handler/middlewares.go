package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/urfave/negroni/v3"
)

func addMiddlewares(next http.Handler) http.Handler {
	return recoveryMiddleware(
		loggingMiddleware(next),
	)
}

func recoveryMiddleware(next http.Handler) http.Handler {
	n := negroni.New()
	recovery := negroni.NewRecovery()
	recovery.Logger = slog.NewLogLogger(slog.Default().Handler(), slog.LevelError)
	n.Use(recovery)
	n.UseHandler(next)
	return n
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug(fmt.Sprintf("Request received: %s %s\n", r.Method, r.URL.String()))

		lrw := negroni.NewResponseWriter(w)
		next.ServeHTTP(lrw, r)

		statusCode := lrw.Status()
		size := lrw.Size()
		slog.Debug(
			fmt.Sprintf(
				"Request handled: %s %s %d %dbytes\n",
				r.Method,
				r.URL.String(),
				statusCode,
				size,
			))
	})
}
