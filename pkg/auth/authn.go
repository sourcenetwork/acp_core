package auth

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

type principalKeyType int

// principalCtxKey is the pkg local key used to set / get principals from ctx
const principalCtxKey principalKeyType = 0

// ExtractPrincipal returns the Principal set in the current ctx
// Return the Anonymous principal if not set
func ExtractPrincipal(ctx context.Context) (types.Principal, error) {
	any := ctx.Value(principalCtxKey)
	switch p := any.(type) {
	case nil:
		return types.AnonymousPrincipal(), nil
	case types.Principal:
		return p, nil
	default:
		return types.Principal{}, ErrInvalidPrincipal
	}
}

// ExtractPrincipalWithType attempts to extract a principal with the given type
// from the current context.
// Return extracted Principal or an error if the current principal does not match the given Kind
func ExtractPrincipalWithType(ctx context.Context, t types.PrincipalKind) (types.Principal, error) {
	principal, err := ExtractPrincipal(ctx)
	if err != nil {
		return types.Principal{}, err
	}
	if principal.Kind == types.PrincipalKind_Anonymous {
		return types.Principal{}, errors.ErrorType_UNAUTHENTICATED
	}
	if principal.Kind != t {
		msg := fmt.Sprintf("principal type wanted %v: got %v", t, principal.Kind)
		return types.Principal{}, errors.Wrap(msg, ErrPrincipalMismatch)

	}
	return principal, nil
}

// InjectPrincipal sets principal to ctx, returns new context
func InjectPrincipal(ctx context.Context, principal types.Principal) context.Context {
	return context.WithValue(ctx, principalCtxKey, principal)
}

func ExtractAuthenticatedPrincipal(ctx context.Context) (types.Principal, error) {
	p, err := ExtractPrincipal(ctx)
	if err != nil {
		return types.Principal{}, err
	}
	if p.Kind == types.PrincipalKind_Anonymous {
		return types.Principal{}, errors.ErrorType_UNAUTHENTICATED
	}
	return p, nil
}
