package utils

import (
	"context"
	"time"

	"cosmossdk.io/log"
)

const spanEventType = "msg_span"

const (
	AttrStart       = "start"
	AttrDuration    = "duration"
	AttrID          = "span_id"
	AttrTx          = "tx_hash"
	AttrMsg         = "msg"
	AttrChainID     = "chain_id"
	AttrBlockHeight = "block_height"
)

type spanCtxKeyType struct{}

var spanCtxKey spanCtxKeyType = struct{}{}

// MsgSpan is a container for execution data generated during the processing of a Msg
// MsgSpan is used as tracing data and must not be relied upon by users, it's an introspection tool only.
// Attributes are not guaranteed to be stable or deterministic
type MsgSpan struct {
	start    time.Time
	duration time.Duration

	tx         string
	id         string
	message    string
	attributes map[string]string
}

func NewSpan(ctx context.Context) *MsgSpan {
	span := &MsgSpan{
		start:      time.Now(),
		attributes: make(map[string]string),
	}
	/*
		hasher := sha256.New()
		hasher.Write(ctx.TxBytes())
		hash := hasher.Sum(nil)
		span.Attr(AttrTx, hex.EncodeToString(hash))
		span.Attr(AttrID, uuid.NewString())
		span.Attr(AttrChainID, ctx.ChainID())
		span.Attr(AttrBlockHeight, fmt.Sprint(ctx.BlockHeight()))
	*/

	return span
}

func (s *MsgSpan) End() {
	s.duration = time.Since(s.start)
}

func (s *MsgSpan) SetMessage(msg string) {
	s.message = msg
}

func (s *MsgSpan) Attr(key, value string) {
	s.attributes[key] = value
}

/*
func (s *MsgSpan) ToEvent() sdk.Event {
	var attrs []sdk.Attribute
	attrs = append(attrs, sdk.NewAttribute(AttrStart, s.start.String()))
	attrs = append(attrs, sdk.NewAttribute(AttrDuration, s.duration.String()))
	if s.message != "" {
		attrs = append(attrs, sdk.NewAttribute(AttrMsg, s.message))
	}

	for key, value := range s.attributes {
		attrs = append(attrs, sdk.NewAttribute(key, value))
	}

	return sdk.NewEvent(spanEventType, attrs...)
}
*/

func (s *MsgSpan) Log(logger log.Logger) {
	var kvs []any
	kvs = append(kvs, AttrStart)
	kvs = append(kvs, s.start.String())
	kvs = append(kvs, AttrDuration)
	kvs = append(kvs, s.duration.String())

	for key, value := range s.attributes {
		kvs = append(kvs, key)
		kvs = append(kvs, value)
	}

	logger.Info("SPAN", kvs...)
}

// WithMsgSpan returns a new Context with an initialized MsgSpan
func WithMsgSpan(ctx context.Context) (context.Context, *MsgSpan) {
	var span *MsgSpan

	spanCtx := ctx.Value(spanCtxKey)
	if spanCtx == nil {
		span = NewSpan(ctx)
		ctx = context.WithValue(ctx, spanCtxKey, span)
	} else {
		span = spanCtx.(*MsgSpan)
	}

	return ctx, span
}

func GetMsgSpan(ctx context.Context) *MsgSpan {
	return ctx.Value(spanCtxKey).(*MsgSpan)
}
