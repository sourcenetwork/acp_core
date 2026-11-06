package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/services"
)

const samplePolicy = `name: folders
resources:
- name: folder
  relations:
  - name: parent
    types: [folder]
  - name: viewer
    types: [actor]
  permissions:
  - name: view
    expr: viewer + parent->view
- name: doc
  relations:
  - name: folder
    types: [folder]
  - name: viewer
    types: [actor]
  permissions:
  - name: view
    expr: viewer + folder->view
`

func post(t *testing.T, h http.Handler, method, body string) *httptest.ResponseRecorder {
	t.Helper()
	return postTo(t, h, method, GeneratePath, body)
}

func postTo(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestGenerate(t *testing.T) {
	if _, err := exec.LookPath("alloy"); err != nil {
		t.Skip("alloy not in PATH")
	}
	h := NewHandler(services.NewPlayground2Service(""))
	body := `{"systemModel":{"policy":` + jsonString(samplePolicy) + `,"constraints":[{"acyclic":{"resource":"folder","relation":"parent"}},{"atMostOne":{"resource":"doc","relation":"folder"}}]},"config":{"count":2}}`
	rec := post(t, h, http.MethodPost, body)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	require.Contains(t, rec.Body.String(), `"worlds"`)
	require.Contains(t, rec.Body.String(), `"relationships"`)
}

func TestInvalidPolicy(t *testing.T) {
	h := NewHandler(services.NewPlayground2Service(""))
	rec := post(t, h, http.MethodPost, `{"systemModel":{"policy":"not: [valid"}}`)
	require.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
	require.Contains(t, rec.Body.String(), `"code":"invalid_argument"`)
	require.Contains(t, rec.Body.String(), `"message"`)
}

func TestBadJSON(t *testing.T) {
	h := NewHandler(services.NewPlayground2Service(""))
	rec := post(t, h, http.MethodPost, `{nope`)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), `"invalid_argument"`)
}

func TestMethodNotAllowed(t *testing.T) {
	h := NewHandler(services.NewPlayground2Service(""))
	rec := post(t, h, http.MethodGet, "")
	require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

const folderWorld = "folder:specs#owner@did:user:alice\nfolder:specs#viewer@did:user:bob\ndoc:roadmap#owner@did:user:carol\ndoc:roadmap#folder@folder:specs\n"

func checkBody(world string) string {
	return `{"scenario":{"systemModel":{"policy":` + jsonString(samplePolicy) + `},"world":{"relationships":` + jsonString(world) +
		`}},"object":{"resource":"doc","id":"roadmap"},"permission":"view","actor":{"id":"did:user:bob"}}`
}

func TestCheck(t *testing.T) {
	h := NewHandler(services.NewPlayground2Service(""))
	rec := postTo(t, h, http.MethodPost, CheckPath, checkBody(folderWorld))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	require.JSONEq(t, `{"allowed":true}`, rec.Body.String())
}

func TestCheckInvalidWorld(t *testing.T) {
	h := NewHandler(services.NewPlayground2Service(""))
	rec := postTo(t, h, http.MethodPost, CheckPath, checkBody("doc:roadmap#viewer@did:user:bob\n"))
	require.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
	require.Contains(t, rec.Body.String(), `"code":"invalid_argument"`)
	require.Contains(t, rec.Body.String(), "object not registered")
}

func TestExplain(t *testing.T) {
	h := NewHandler(services.NewPlayground2Service(""))
	rec := postTo(t, h, http.MethodPost, ExplainPath, checkBody(folderWorld))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var out struct {
		Allowed bool `json:"allowed"`
		Graph   struct {
			Nodes []json.RawMessage `json:"nodes"`
			Edges []json.RawMessage `json:"edges"`
		} `json:"graph"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &out))
	require.True(t, out.Allowed)
	require.NotEmpty(t, out.Graph.Nodes)
	require.NotEmpty(t, out.Graph.Edges)
}

func TestCheckMethodNotAllowed(t *testing.T) {
	h := NewHandler(services.NewPlayground2Service(""))
	require.Equal(t, http.StatusMethodNotAllowed, postTo(t, h, http.MethodGet, CheckPath, "").Code)
}

func jsonString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func TestCompile(t *testing.T) {
	h := NewHandler(services.NewPlayground2Service(""))
	body := `{"policy":` + jsonString(samplePolicy) + `,"marshalType":"YAML","module":"policy","constraints":[{"acyclic":{"resource":"folder","relation":"parent"}}]}`
	rec := postTo(t, h, http.MethodPost, CompilePath, body)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	for _, name := range []string{"policy.als", "policy_management.als", "policy_assumptions.als"} {
		require.Contains(t, rec.Body.String(), `"name":"`+name+`"`)
	}
	require.Contains(t, rec.Body.String(), "fun doc_view")
}
