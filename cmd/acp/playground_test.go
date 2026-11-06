package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/internal/workspace"
)

func newTestPlayground(t *testing.T, backend string) (*httptest.Server, string) {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "policy.yaml"), []byte("name: notes\n"), 0644))
	target, err := url.Parse(backend)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ui := fstest.MapFS{"index.html": {Data: []byte("<html>playground</html>")}}
	srv := httptest.NewServer(newPlaygroundHandler(ctx, dir, target, ui, 10*time.Millisecond))
	t.Cleanup(srv.Close)
	return srv, dir
}

func Test_Playground_ServesWorkspaceAndUI(t *testing.T) {
	srv, dir := newTestPlayground(t, "http://localhost:1")

	resp, err := http.Get(srv.URL + workspacePath)
	require.NoError(t, err)
	defer resp.Body.Close()
	var files workspace.Files
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&files))
	require.Equal(t, "name: notes\n", files.Policy)
	require.Equal(t, dir, files.Dir)

	resp, err = http.Get(srv.URL + "/")
	require.NoError(t, err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), "playground")
}

func Test_Playground_ProxiesTheService(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		io.WriteString(w, `{"path":"`+r.URL.Path+`","body":`+string(data)+`}`)
	}))
	defer backend.Close()
	srv, _ := newTestPlayground(t, backend.URL)

	resp, err := http.Post(srv.URL+"/sourcenetwork.acp_core.Playground2/Verify", "application/json", strings.NewReader(`{"x":1}`))
	require.NoError(t, err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, `{"path":"/sourcenetwork.acp_core.Playground2/Verify","body":{"x":1}}`, string(body))
}

func Test_Playground_ServiceDown(t *testing.T) {
	srv, _ := newTestPlayground(t, "http://localhost:1")
	resp, err := http.Post(srv.URL+"/sourcenetwork.acp_core.Playground2/Verify", "application/json", strings.NewReader(`{}`))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadGateway, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), `"code":"unavailable"`)
	require.Contains(t, string(body), "can't reach the acp service")
}

func Test_Playground_EventsOnChange(t *testing.T) {
	srv, dir := newTestPlayground(t, "http://localhost:1")
	resp, err := http.Get(srv.URL + eventsPath)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))
	events := bufio.NewReader(resp.Body)

	next := func() string {
		t.Helper()
		var event string
		for {
			line, err := events.ReadString('\n')
			require.NoError(t, err)
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "event: ") {
				event = strings.TrimPrefix(line, "event: ")
			}
			if line == "" && event != "" {
				return event
			}
		}
	}
	require.Equal(t, "ready", next())

	require.NoError(t, os.MkdirAll(filepath.Join(dir, "worlds"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "worlds", "a.yaml"), []byte("relationships: x\n"), 0644))
	require.Equal(t, "changed", next())
}
