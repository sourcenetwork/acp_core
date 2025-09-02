// package decorator defines types which are used to decorate request handlers
package decorator

import "context"

// TypedHandler models a strongly typed request handler
type TypedHandler[Req, Resp any] func(ctx context.Context, req Req) (Resp, error)

// Handler models a request handler
type Handler func(ctx context.Context, req any) (any, error)

// Decorator is function that takes a Handler and returns a new Handler
type Decorator func(Handler) Handler

// ToTypedHandler takes a generic params and converts a Handler to a TypedHandler
func ToTypedHandler[Req, Resp any](h Handler) TypedHandler[Req, Resp] {
	return func(ctx context.Context, req Req) (Resp, error) {
		resp, err := h(ctx, req)
		return resp.(Resp), err
	}
}

// ToAnyHandler removes the strong typed guarantees of a TypedHandler
func ToAnyHandler[Req, Resp any](h TypedHandler[Req, Resp]) Handler {
	return func(ctx context.Context, req any) (any, error) {
		resp, err := h(ctx, req.(Req))
		return resp, err
	}
}

// Decorate applies decorators to h
func Decorate(h Handler, decorators ...Decorator) Handler {
	for _, dec := range decorators {
		h = dec(h)
	}
	return h
}

// Chain takes a sequence of Decorators and returns one composite Decorator
func Chain(decorators ...Decorator) Decorator {
	return func(h Handler) Handler {
		for _, dec := range decorators {
			h = dec(h)
		}
		return h
	}
}

// DecorateTypedHandler decorates a TypedHandler
func DecorateTypedHandler[Req, Resp any](h TypedHandler[Req, Resp], dec Decorator) TypedHandler[Req, Resp] {
	return ToTypedHandler[Req, Resp](Decorate(ToAnyHandler(h), dec))
}
