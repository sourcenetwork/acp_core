package playground

import (
	"context"

	"github.com/sourcenetwork/acp_core/pkg/playground"
)

func NewSandbox(ctx context.Context, in *playground.NewSandboxRequest) (*playground.NewSandboxResponse, error) {
	return nil, nil
}

// ListSandboxes returns the list of sandboxes created in the Playgruond
func ListSandboxes(ctx context.Context, in *playground.ListSandboxesRequest) (*playground.ListSandboxesResponse, error) {
	return nil, nil
}

// SetState updates the state of a Sandbox environment with the newly provided data
func SetState(ctx context.Context, in *playground.SetStateRequest) (*playground.SetStateResponse, error) {
	return nil, nil
}

// GetCatalogue returns the catalogue (index/LUT) of defined entities within a sandbox
func GetCatalogue(ctx context.Context, in *playground.GetCatalogueRequest) (*playground.GetCatalogueResponse, error) {
	return nil, nil
}

// VerifyTheorems executes the defined test suite for a sandbox
func VerifyTheorems(ctx context.Context, in *playground.VerifyTheoremsRequest) (*playground.VerifyTheoremsResponse, error) {
	return nil, nil
}
