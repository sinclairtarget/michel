package server

import (
	"fmt"
	"log/slog"
	"math"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/sinclairtarget/michel/internal/util"
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

// Starts goroutine for handling fsnotify events.
//
// We debounce events to make sure we don't prematurely handle a change to a
// file (e.g. by reacting to the first of several write events).
func (w *watcher) start() error {
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

				wrappedEv, err := w.handle(ev)
				if err != nil {
					slog.Error("error handling fsnotify event", "error", err)
					close(w.events)
					return
				}

				if wrappedEv != nil {
					timer.Stop()
					timer = time.AfterFunc(debounceMs, func() {
						w.events <- *wrappedEv
					})
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					timer.Stop()
					close(w.events)
					return
				}
				slog.Error("fsnotify error", "error", err)
			case <-w.done:
				timer.Stop()
				close(w.events)
				return
			}
		}
	}()

	// Set up initial watch list
	for _, dir := range w.dirs {
		seq, finish := util.WalkDirs(dir)

		// Dir will be a subdir of itself, so it gets added too
		for subdir := range seq {
			err := w.add(subdir)
			if err != nil {
				return err
			}
		}

		err := finish()
		if err != nil {
			return err
		}
	}

	return nil
}

// Logic for handling fsnotify events.
//
// Returns a nil event if the event should be ignored.
func (w *watcher) handle(ev fsnotify.Event) (*event, error) {
	slog.Debug("fsnotify event", "event", ev)
	if ev.Has(fsnotify.Chmod) {
		return nil, nil
	}

	if ev.Has(fsnotify.Create) {
		isDir, err := util.IsDir(ev.Name)
		if err != nil {
			return nil, err
		}

		if isDir {
			// A new directory! We want to watch it too
			err := w.add(ev.Name)
			if err != nil {
				return nil, err
			}
		}
	}

	return &event{
		path: ev.Name,
	}, nil
}

func (w *watcher) add(path string) error {
	// Not necessary for fsnotify, but it's nice to have the absolute path for
	// logging.
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	slog.Debug("adding path to watch list", "filepath", absPath)

	err = w.watcher.Add(absPath)
	return err
}

func (w *watcher) close() {
	close(w.done)
	w.watcher.Close()
}
