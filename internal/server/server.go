/*
* Local HTTP server that can serve blog pages.
*
* Watches site directories for file changes and triggers a full rebuild on any
* change.
*
* TODO: Make build incremental.
 */
package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/sinclairtarget/michel/internal/build"
)

func Run(logger *slog.Logger, basePath string, port int) error {
	watcher := newWatcher(
		build.ContentDir,
		build.LayoutsDir,
		build.PartialsDir,
		build.PagesDir,
	)
	defer watcher.close_()

	// Goroutine to watch for changes.
	// Triggers a full build for any change.
	go func() {
		for event := range watcher.events {
			logger.Debug("got file modified event", "path", event.path)
			rebuild(logger)
		}

		logger.Debug("goroutine exiting; watch events channel closed")
	}()

	err := watcher.start(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start file watcher: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Starting server on port %d...\n", port)

	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(
		addr,
		logMiddleware(logger, http.FileServer(http.Dir(basePath)).ServeHTTP),
	)
}

func rebuild(logger *slog.Logger) {
	start := time.Now()
	err := build.Build(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during build: %v\n", err)
	}

	elapsed := time.Now().Sub(start)
	fmt.Printf("Site rebuilt in %dms.\n", elapsed.Milliseconds())
}

func logMiddleware(logger *slog.Logger, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(logger, r)
		f(w, r)
	}
}

func logRequest(logger *slog.Logger, r *http.Request) {
	logger.Info("http request", "method", r.Method, "url", r.URL)
}
