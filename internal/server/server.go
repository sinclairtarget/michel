package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

// bim
func logMiddleware(logger *slog.Logger, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("GET", "url", r.URL)
		f(w, r)
	}
}

func Run(logger *slog.Logger, basePath string, port int) error {
	fmt.Printf("Starting server on port %d...\n", port)

	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(
		addr,
		logMiddleware(logger, http.FileServer(http.Dir(basePath)).ServeHTTP),
	)
}
