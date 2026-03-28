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

func Run(bind string, port int, outdir string) error {
	watcher := newWatcher(
		build.ContentDir,
		build.LayoutsDir,
		build.PartialsDir,
		build.SiteDir,
	)
	defer watcher.close()

	// Goroutine to watch for changes.
	// Triggers a full build for any change.
	go func() {
		for event := range watcher.events {
			slog.Debug("got file modified event", "path", event.path)
			rebuild(outdir)
		}

		slog.Debug("goroutine exiting; watch events channel closed")
	}()

	err := watcher.start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start file watcher: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Starting server on port %d...\n", port)

	addr := fmt.Sprintf("%s:%d", bind, port)
	return http.ListenAndServe(
		addr,
		logMiddleware(http.FileServer(http.Dir(outdir)).ServeHTTP),
	)
}

func rebuild(outdir string) {
	start := time.Now()
	err := build.Build(outdir)
	if err != nil {
		build.PrintBuildError(err)
	}

	elapsed := time.Now().Sub(start)
	fmt.Printf("Site rebuilt in %dms.\n", elapsed.Milliseconds())
}

func logMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		f(w, r)
	}
}

func logRequest(r *http.Request) {
	slog.Info("http request", "method", r.Method, "url", r.URL)
}
