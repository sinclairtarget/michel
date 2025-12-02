package watch

import (
	"fmt"
	"log/slog"
	"math"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

var debounceMs time.Duration = 1000 * time.Millisecond

type Event struct {
	Path string
}

type Watcher struct {
	Path    string
	Events  chan Event
	done    chan struct{}
	watcher *fsnotify.Watcher
}

func NewWatcher(path string) Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Sprintf("failed to create fsnotify watcher: %v", err))
	}

	return Watcher{
		Path:    path,
		Events:  make(chan Event),
		done:    make(chan struct{}),
		watcher: watcher,
	}
}

func (w *Watcher) Start(logger *slog.Logger) error {
	go func() {
		timer := time.AfterFunc(math.MaxInt64, func() {})

		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					timer.Stop()
					close(w.Events)
					return
				}

				logger.Debug("fsnotify event", "event", event)
				if !event.Has(fsnotify.Chmod) {
					timer.Stop()
					timer = time.AfterFunc(debounceMs, func() {
						w.Events <- Event{
							Path: event.Name,
						}
					})
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					timer.Stop()
					close(w.Events)
					return
				}
				logger.Error("fsnotify error", "error", err)
			case <-w.done:
				timer.Stop()
				close(w.Events)
				return
			}
		}
	}()

	return w.add(logger, w.Path)
}

func (w *Watcher) add(logger *slog.Logger, path string) error {
	absPath, err := filepath.Abs(w.Path)
	if err != nil {
		return err
	}

	logger.Debug("adding fsnotify path", "filepath", absPath)

	err = w.watcher.Add(w.Path)
	return err
}

func (w *Watcher) Close() {
	close(w.done)
	w.watcher.Close()
}
