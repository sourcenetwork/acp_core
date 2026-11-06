package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cosmos/gogoproto/jsonpb"
	"github.com/cosmos/gogoproto/proto"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

// client calls the Playground2 service over the Connect unary JSON protocol
type client struct {
	server string
	http   *http.Client
}

func newClient(server string) *client {
	// verifying runs alloy once per contract, which can take a while
	return &client{server: strings.TrimSuffix(server, "/"), http: &http.Client{Timeout: 5 * time.Minute}}
}

func (c *client) Verify(ctx context.Context, req *types.VerifyRequest) (*types.VerifyResponse, error) {
	resp := &types.VerifyResponse{}
	return resp, c.call(ctx, "Verify", req, resp)
}

func (c *client) Generate(ctx context.Context, req *types.GenerateRequest) (*types.GenerateResponse, error) {
	resp := &types.GenerateResponse{}
	return resp, c.call(ctx, "Generate", req, resp)
}

func (c *client) Check(ctx context.Context, req *types.CheckRequest) (*types.CheckResponse, error) {
	resp := &types.CheckResponse{}
	return resp, c.call(ctx, "Check", req, resp)
}

func (c *client) Explain(ctx context.Context, req *types.ExplainRequest) (*types.ExplainResponse, error) {
	resp := &types.ExplainResponse{}
	return resp, c.call(ctx, "Explain", req, resp)
}

func (c *client) Compile(ctx context.Context, req *types.CompileRequest) (*types.CompileResponse, error) {
	resp := &types.CompileResponse{}
	return resp, c.call(ctx, "Compile", req, resp)
}

func (c *client) call(ctx context.Context, method string, req, resp proto.Message) error {
	body, err := (&jsonpb.Marshaler{}).MarshalToString(req)
	if err != nil {
		return fmt.Errorf("encoding request: %w", err)
	}
	url := c.server + "/sourcenetwork.acp_core.Playground2/" + method
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := c.http.Do(httpReq)
	if err != nil {
		return fmt.Errorf("can't reach the acp service at %v (set -server or ACP_SERVER): %w", c.server, err)
	}
	defer httpResp.Body.Close()
	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		var connectErr struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if json.Unmarshal(data, &connectErr) == nil && connectErr.Message != "" {
			if connectErr.Code == "invalid_argument" {
				return fmt.Errorf("%v", connectErr.Message)
			}
			return fmt.Errorf("acp service error (%v): %v", connectErr.Code, connectErr.Message)
		}
		return fmt.Errorf("acp service error: %v", httpResp.Status)
	}
	if err := (&jsonpb.Unmarshaler{AllowUnknownFields: true}).Unmarshal(bytes.NewReader(data), resp); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}
