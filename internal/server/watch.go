package server

import (
	"fmt"
	"log/slog"
	"math"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

var debounceMs time.Duration = 1000 * time.Millisecond

type event struct {
	path string
}

type watcher struct {
	dirs    []string
	events  chan event
	done    chan struct{}
	watcher *fsnotify.Watcher
}

func newWatcher(dirs ...string) watcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Sprintf("failed to create fsnotify watcher: %v", err))
	}

	return watcher{
		dirs:    dirs,
		events:  make(chan event),
		done:    make(chan struct{}),
		watcher: w,
	}
}

func (w *watcher) start(logger *slog.Logger) error {
	go func() {
		timer := time.AfterFunc(math.MaxInt64, func() {})

		for {
			select {
			case ev, ok := <-w.watcher.Events:
				if !ok {
					timer.Stop()
					close(w.events)
					return
				}

				logger.Debug("fsnotify event", "event", ev)
				if !ev.Has(fsnotify.Chmod) {
					timer.Stop()
					timer = time.AfterFunc(debounceMs, func() {
						w.events <- event{
							path: ev.Name,
						}
					})
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					timer.Stop()
					close(w.events)
					return
				}
				logger.Error("fsnotify error", "error", err)
			case <-w.done:
				timer.Stop()
				close(w.events)
				return
			}
		}
	}()

	for _, dir := range w.dirs {
		err := w.add(logger, dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *watcher) add(logger *slog.Logger, path string) error {
	// Not necessary for fsnotify, but it's nice to have the absolute path for
	// logging.
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	logger.Debug("adding fsnotify path", "filepath", absPath)

	err = w.watcher.Add(absPath)
	return err
}

func (w *watcher) close_() {
	close(w.done)
	w.watcher.Close()
}
