package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

func Run(logger *slog.Logger, basePath string, port int) error {
	fmt.Printf("Starting server on port %d...\n", port)

	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(
		addr,
		logMiddleware(logger, http.FileServer(http.Dir(basePath)).ServeHTTP),
	)
}

func logMiddleware(logger *slog.Logger, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(logger, r)
		f(w, r)
	}
}

func logRequest(logger *slog.Logger, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL)
	logger.Debug("http request", "method", r.Method, "url", r.URL)
}
