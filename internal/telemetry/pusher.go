package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/version"
)

const jsonContentType = "application/json"

// ErrorPayload is a copy of backend `PublishErrorRequest“ message
type ErrorPayload struct {
	Message           string             `json:"message"`
	StackTrace        string             `json:"stack_trace"`
	State             *types.SandboxData `json:"state"`
	PlaygroundVersion string             `json:"playground_version"`
}

// NetPlaygroundBackendErrorClient returns a client to publish errors
func NewPlaygroundBackendErrorClient(url string) *PlaygroundBackendErrorClient {
	return &PlaygroundBackendErrorClient{
		url: url,
	}
}

// PlaygroundBackendErrorClient models a client which communicates with an instance of
// the PlaygroundBackend service in order to publish errors
type PlaygroundBackendErrorClient struct {
	url string
}

// PublishError submits an error to the Backend API
func (c *PlaygroundBackendErrorClient) PushError(ctx context.Context, data *types.SandboxData, err error) error {
	// looks for a stack trace attribute in the error
	trace := ""
	if acpErr, ok := err.(*errors.Error); ok {
		t, ok := acpErr.Metadata["stack_trace"]
		if ok {
			trace = t
		}
	}
	payload := ErrorPayload{
		Message:           err.Error(),
		StackTrace:        trace,
		State:             data,
		PlaygroundVersion: version.Commit,
	}
	bytes, err := json.Marshal(&payload)
	if err != nil {
		return fmt.Errorf("marshaling payload: %v", err)
	}
	reader := strings.NewReader(string(bytes))
	resp, err := http.Post(c.url, jsonContentType, reader)
	if err != nil {
		return fmt.Errorf("error publishing error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error publishing error: %v", err)
	}
	return nil
}
