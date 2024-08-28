//go:build js

package playground_js

import (
	"syscall/js"

	"github.com/sourcenetwork/acp_core/pkg/errors"
)

// asyncFn takes a Go Function which returns JS Values and turns it into a JS function which returns a promise.
// First a Promise executor function is created, inside of which a worker go routine is created and dispatched.
// The executor returns immediately and the returned Promise is returned as the result of the function
//
// Inside the worker go routine, the Handler is invoked and according to the results
// either invokes the reject if the error wasn't nil or resolves the promise with the result
func asyncFn(f func(js.Value, []js.Value) (js.Value, error)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		var promiseExecutor js.Func
		promiseExecutor = js.FuncOf(func(_ js.Value, promiseArgs []js.Value) any {
			go func() {
				defer promiseExecutor.Release()
				resolve := promiseArgs[0]
				reject := promiseArgs[1]

				result, err := f(this, args)
				if err != nil {
					jsErr := errors.NewJSError(err)
					reject.Invoke(jsErr)
					return
				}
				resolve.Invoke(result)
			}()
			return js.Undefined()
		})
		return js.Global().Get("Promise").New(promiseExecutor)
	})
}
