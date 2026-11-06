package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sourcenetwork/acp_core/internal/workspace"
)

const (
	// workspacePath serves the workspace's files, as workspace.Files
	workspacePath = "/api/workspace"
	// eventsPath streams server-sent events: `ready` on connect, then `changed` whenever a file changes,
	// both with data {"version": "..."}
	eventsPath = "/api/workspace/events"
	// servicePrefix is proxied to the acp service
	servicePrefix = "/sourcenetwork.acp_core.Playground2/"

	pollInterval      = 500 * time.Millisecond
	keepaliveInterval = 25 * time.Second
)

// playground serves the playground UI for the workspace in dir, proxying API calls to the acp service
func playground(server, dir, addr, uiDir string, open bool) {
	target, err := url.Parse(server)
	if err != nil || target.Scheme == "" || target.Host == "" {
		fatal(fmt.Sprintf("invalid server address %q: expected eg. http://localhost:8090", server))
	}
	if _, err := workspace.ReadFiles(dir); err != nil {
		fatal(err)
	}
	ui := builtinUI()
	if uiDir != "" {
		ui = os.DirFS(uiDir)
	}

	ctx := context.Background()
	h := newPlaygroundHandler(ctx, dir, target, ui, pollInterval)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fatal(err)
	}
	pageURL := "http://" + listener.Addr().String()
	log.Printf("playground for %v at %v (acp service: %v)", dir, pageURL, server)
	if open {
		openBrowser(pageURL)
	}
	fatal(http.Serve(listener, h))
}

// newPlaygroundHandler watches the workspace in dir until ctx is done, polling every interval
func newPlaygroundHandler(ctx context.Context, dir string, target *url.URL, ui fs.FS, interval time.Duration) http.Handler {
	h := &hub{subs: make(map[chan string]struct{})}
	go h.watch(ctx, dir, interval)

	mux := http.NewServeMux()
	mux.HandleFunc("GET "+workspacePath, func(w http.ResponseWriter, r *http.Request) {
		files, err := workspace.ReadFiles(dir)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"code": "internal", "message": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, files)
	})
	mux.HandleFunc("GET "+eventsPath, h.serveEvents)
	mux.Handle(servicePrefix, &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			writeJSON(w, http.StatusBadGateway, map[string]string{
				"code":    "unavailable",
				"message": fmt.Sprintf("can't reach the acp service at %v: %v", target, err),
			})
		},
	})
	mux.Handle("/", http.FileServer(http.FS(ui)))
	return mux
}

// hub tells subscribers about workspace changes
type hub struct {
	mu      sync.Mutex
	version string
	subs    map[chan string]struct{}
}

func (h *hub) watch(ctx context.Context, dir string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		// a workspace which fails to load (eg. policy.yaml mid-rename) is a version too: clients refetch and show the error
		version := "unreadable"
		if files, err := workspace.ReadFiles(dir); err == nil {
			version = files.Version
		}
		h.mu.Lock()
		if version != h.version {
			h.version = version
			for sub := range h.subs {
				select {
				case sub <- version:
				default: // the subscriber has a change pending already: it'll fetch the latest files anyway
				}
			}
		}
		h.mu.Unlock()

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (h *hub) serveEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	sub := make(chan string, 1)
	h.mu.Lock()
	h.subs[sub] = struct{}{}
	version := h.version
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		delete(h.subs, sub)
		h.mu.Unlock()
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	send := func(event, version string) {
		data, _ := json.Marshal(map[string]string{"version": version})
		fmt.Fprintf(w, "event: %v\ndata: %s\n\n", event, data)
		flusher.Flush()
	}
	send("ready", version)

	keepalive := time.NewTicker(keepaliveInterval)
	defer keepalive.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case version := <-sub:
			send("changed", version)
		case <-keepalive.C:
			fmt.Fprint(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func openBrowser(pageURL string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", pageURL)
	case "darwin":
		cmd = exec.Command("open", pageURL)
	default:
		return
	}
	if err := cmd.Start(); err != nil {
		log.Printf("couldn't open a browser (%v): open %v", strings.TrimSpace(err.Error()), pageURL)
	}
}
