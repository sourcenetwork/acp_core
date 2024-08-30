// package playground_js defines a wrapper for a PlaygroundService which receives inputs and return outputs to the JS runtime
//
// The primary type of playgrond_js is PlaygroundServiceProxy, which wraps a PlaygroundService implementation
// that operates on the JS domain.
// In the JS domain, this type can be constructed using the `PlaygroundConstructor` JS function,
// the returned value is a JS object which satisfies the PlaygroundService TypeScritpt interface
// generated from the protobuff descriptors.
//
// Note: the returned object satisfies the Service interface but it also includes an additional method,
// `close`, which should be called if the JS Service is going to be discarded.
// `close` frees resources allocated in the Go domain.
//
// A small note about error handling is required.
// Due to the gap between JS exceptions and Go lang explicit error values, error handling
// for the wrapper is done through a Promise system.
// Every method in the service implementation returns a Promise which will be resolved with an error
// in the event of a Go land error.
// The JS Errors names will match the ErrorTypes defined in pkg/errors.
package playground_js
