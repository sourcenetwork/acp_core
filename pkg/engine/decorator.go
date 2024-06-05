package engine

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/utils"
)

var _ Decorator = (*MsgSpanDecorator)(nil)

// Decorator models a type which acts before or after the execution of a MsgServer
type Decorator interface {
	// Name returns an identifier for the Decorator
	Name() string

	// Pre is called before the msg server processes the Msg
	// The Msg is immutable but the Decorator may update the context
	Pre(ctx context.Context, msg any) (context.Context, error)

	// Post is called after the Msg Server processes the Msg
	// The response msg is immutable but the Decorator may update the context
	Post(ctx context.Context, response any, handlerErr error) (context.Context, error)
}

func NewMsgSpanDecorator(log bool, emitEvent bool) *MsgSpanDecorator {
	return &MsgSpanDecorator{
		log:       log,
		emitEvent: emitEvent,
	}
}

type MsgSpanDecorator struct {
	emitEvent bool
	log       bool
}

func (h *MsgSpanDecorator) Name() string {
	return "MsgSpanDecorator"
}

func (h *MsgSpanDecorator) Pre(ctx context.Context, msg any) (context.Context, error) {
	ctx, _ = utils.WithMsgSpan(ctx)
	return ctx, nil
}

func (h *MsgSpanDecorator) Post(ctx context.Context, resp any, err error) (context.Context, error) {
	span := utils.GetMsgSpan(ctx)
	span.End()

	if h.log {
		//logger := .Logger()
		//span.Log(logger)
	}

	if h.emitEvent {
		//event := span.ToEvent()
		//sdkCtx.EventManager().EmitEvent(event)
	}

	return ctx, nil
}

type Handler[Req, Resp any] func(context.Context, Req) (Resp, error)

func runPre[Req any](ctx context.Context, decorators []Decorator, msg Req) (context.Context, error) {
	var err error
	for _, hook := range decorators {
		ctx, err = hook.Pre(ctx, msg)
		if err != nil {
			return ctx, fmt.Errorf("hook %v: pre exec: %w", hook.Name(), err)
		}
	}
	return ctx, nil
}

func runPost[Resp any](ctx context.Context, decorators []Decorator, resp Resp, handlerErr error) error {
	var hookErr error
	for _, hook := range decorators {
		ctx, hookErr = hook.Post(ctx, resp, handlerErr)
		if hookErr != nil {
			return fmt.Errorf("hook %v: post exec: %w", hook.Name(), hookErr)
		}
	}
	return nil
}

func applyMiddleware[Req, Resp any](ctx context.Context, handler Handler[Req, Resp], decorators []Decorator, req Req) (Resp, error) {
	var zero Resp
	ctx, err := runPre(ctx, decorators, req)
	if err != nil {
		return zero, err
	}

	resp, err := handler(ctx, req)

	postErr := runPost(ctx, decorators, resp, err)

	if err != nil {
		return zero, err
	}
	if postErr != nil {
		return zero, postErr
	}

	return resp, nil
}
