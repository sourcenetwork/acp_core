package main

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"
	"io"
	"net/http"

	"github.com/cosmos/gogoproto/jsonpb"
	"github.com/cosmos/gogoproto/proto"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// Connect unary JSON routes for Playground2.
const (
	GeneratePath = "/sourcenetwork.acp_core.Playground2/Generate"
	CheckPath    = "/sourcenetwork.acp_core.Playground2/Check"
	ExplainPath  = "/sourcenetwork.acp_core.Playground2/Explain"
	VerifyPath   = "/sourcenetwork.acp_core.Playground2/Verify"
	CompilePath  = "/sourcenetwork.acp_core.Playground2/Compile"
)

const maxBodyBytes = 4 << 20

// NewHandler serves Playground2 over the Connect unary JSON protocol.
func NewHandler(svc types.Playground2Server) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(GeneratePath, unary(svc.Generate))
	mux.Handle(CheckPath, unary(svc.Check))
	mux.Handle(ExplainPath, unary(svc.Explain))
	mux.Handle(VerifyPath, unary(svc.Verify))
	mux.Handle(CompilePath, unary(svc.Compile))
	return mux
}

// unary adapts a unary service method to a Connect JSON POST handler.
// Req and Resp are the pointer types of the generated messages.
func unary[Req, Resp any, PReq interface {
	*Req
	proto.Message
}, PResp interface {
	*Resp
	proto.Message
}](call func(context.Context, PReq) (PResp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeError(w, http.StatusMethodNotAllowed, "unimplemented", "method not allowed, use POST")
			return
		}
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodyBytes))
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_argument", "reading request body: "+err.Error())
			return
		}
		req := PReq(new(Req))
		if err := (&jsonpb.Unmarshaler{}).Unmarshal(bytes.NewReader(body), req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_argument", "invalid request JSON: "+err.Error())
			return
		}
		resp, err := call(r.Context(), req)
		if err != nil {
			if stderrors.Is(err, errors.ErrorType_BAD_INPUT) {
				writeError(w, http.StatusBadRequest, "invalid_argument", err.Error())
				return
			}
			writeError(w, http.StatusInternalServerError, "internal", err.Error())
			return
		}
		out, err := (&jsonpb.Marshaler{EmitDefaults: true}).MarshalToString(resp)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal", "encoding response: "+err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, out)
	}
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"code": code, "message": msg})
}
