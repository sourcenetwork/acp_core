//go:build js

package playground_js

import (
	"syscall/js"

	"github.com/cosmos/gogoproto/jsonpb"
	"github.com/cosmos/gogoproto/proto"
	"github.com/sourcenetwork/acp_core/pkg/errors"
)

// unmarshalArgs extracts the first arg from the slice of args,
// dumps it is a JSON string and attempts to unmarshal the JSON string into the proto container
// returns an error if args has more than one element or if the unmarshaling failed.
func unmarsahlArgs(container proto.Message, args []js.Value) error {
	count := len(args)
	if count != 1 {
		return newInvalidArgsErr(count)
	}

	arg := args[0]

	reqJson := js.Global().Get("JSON").Call("stringify", arg)
	if reqJson.Type() != js.TypeString {
		panic("expected string")
	}

	err := jsonpb.UnmarshalString(reqJson.String(), container)
	if err != nil {
		return errors.NewFromCause("could not unmarshal serialized req obj", err, errors.ErrorType_BAD_INPUT)
	}

	return nil
}

// toJSObject maps a proto value into a JS value.
// The proto value is dumped as a JSON string and loaded in the JS runtime
// using the JSON.parse method.
func toJSObject[T proto.Message](val T) (js.Value, error) {
	marshaler := jsonpb.Marshaler{}
	valStr, err := marshaler.MarshalToString(val)
	if err != nil {
		return js.Value{}, err
	}
	jsVal := js.Global().Get("JSON").Call("parse", valStr)
	return jsVal, nil
}
