package test

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertResults is a general helper for computations which return a result and an error
// if wantErr is specified, it asserts that gotErr is a match
// otherwise, it asserts got and want are a match and that gotErr is nil.
func AssertResults(ctx *TestCtx, got, want any, gotErr, wantErr error) {
	if wantErr != nil {
		require.NotNil(ctx.T, gotErr, "expected an error but got none")
		if errors.Is(gotErr, wantErr) {
			assert.ErrorIs(ctx.T, gotErr, wantErr)
		} else {
			// Errors returned from SDK operations (RPC communication to a SourceHub node)
			// no longer have the original errors wrapped, therefore we compare a string as fallback strat.
			gotErrStr := gotErr.Error()
			wantErrStr := wantErr.Error()
			assert.Contains(ctx.T, gotErrStr, wantErrStr)
		}
	} else {
		assert.NoError(ctx.T, gotErr)
	}
	if !isNil(want) {
		assert.Equal(ctx.T, want, got)
	}
}

func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	isNilableKind := containsKind(
		[]reflect.Kind{
			reflect.Chan, reflect.Func,
			reflect.Interface, reflect.Map,
			reflect.Ptr, reflect.Slice, reflect.UnsafePointer},
		kind)

	if isNilableKind && value.IsNil() {
		return true
	}

	return false
}

// containsKind checks if a specified kind in the slice of kinds.
func containsKind(kinds []reflect.Kind, kind reflect.Kind) bool {
	for i := 0; i < len(kinds); i++ {
		if kind == kinds[i] {
			return true
		}
	}

	return false
}

func MustProtoToJson(msg proto.Message) string {
	bytes, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
