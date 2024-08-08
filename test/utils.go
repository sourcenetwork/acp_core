package test

import (
	"errors"
	"reflect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
