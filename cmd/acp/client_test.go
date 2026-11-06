package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

func Test_Client_SendsJSONAndDecodesResponses(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/sourcenetwork.acp_core.Playground2/Verify", r.URL.Path)
		data, _ := io.ReadAll(r.Body)
		body = string(data)
		io.WriteString(w, `{"scope":3,"results":[{"contract":"c","status":"VIOLATED","properties":[{"status":"VIOLATED","world":{"relationships":"x"}}],"newField":1}]}`)
	}))
	defer srv.Close()

	resp, err := newClient(srv.URL+"/").Verify(context.Background(), &types.VerifyRequest{
		SystemModel: &types.SystemModel{Requirements: []*types.Requirement{{Name: "r", Contracts: []*types.Contract{
			{Name: "c", Properties: []*types.Property{{Type: types.PropertyType_PRED}}},
		}}}},
	})
	require.NoError(t, err)
	require.Contains(t, body, `"type":"PRED"`)
	require.Equal(t, types.VerifyStatus_VIOLATED, resp.Results[0].Status)
	require.Equal(t, "x", resp.Results[0].Properties[0].World.Relationships)
}

func Test_Client_Errors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "Verify") {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, `{"code":"invalid_argument","message":"constraint references undefined relation doc.x"}`)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"code":"internal","message":"running alloy: not found"}`)
	}))
	defer srv.Close()
	c := newClient(srv.URL)

	_, err := c.Verify(context.Background(), &types.VerifyRequest{})
	require.EqualError(t, err, "constraint references undefined relation doc.x")
	_, err = c.Generate(context.Background(), &types.GenerateRequest{})
	require.EqualError(t, err, "acp service error (internal): running alloy: not found")
}
