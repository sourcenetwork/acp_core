package auth

import (
	"context"
	"errors"
	"fmt"

	didpkg "github.com/sourcenetwork/acp_core/pkg/did"
)

var ErrInvalidPrincipal = errors.New("invalid principal")
var ErrUnauthenticatd = errors.New("unauthenticated: principal must be authenticated")
var ErrPrincipalMismatch = errors.New("mismatched principal")

type principalKeyType int

var principalCtxKey principalKeyType = 0

const anonymousIdentifier = "anonymous"
const rootIdentifier = "root"

type PrincipalType int

const (
	DID PrincipalType = iota
	Anonymous
	Root
)

func (t PrincipalType) String() string {
	switch t {
	case DID:
		return "DID"
	case Anonymous:
		return "Anonymous"
	case Root:
		return "Root"
	default:
		return ""
	}
}

type AuthProvider interface {
	VerifyActor(ctx context.Context, user string) error
}

type UserRegistryProvider struct {
	users map[string]struct{}
}

func (p *UserRegistryProvider) VerifyActor(ctx context.Context, user string) error {
	_, ok := p.users[user]
	if !ok {
		return fmt.Errorf("user %v: not found", user)
	}

	return nil
}

func NewUserRegistryProvider(users ...string) *UserRegistryProvider {
	usersMap := make(map[string]struct{})
	for _, user := range users {
		usersMap[user] = struct{}{}
	}

	return &UserRegistryProvider{
		users: usersMap,
	}
}

var _ Principal = (*didPrincipal)(nil)

type Principal interface {
	Identifier() string
	GetType() PrincipalType
	IsAnonymous() bool
	IsRoot() bool
}

type didPrincipal string

func (p didPrincipal) Identifier() string     { return string(p) }
func (p didPrincipal) IsAnonymous() bool      { return false }
func (p didPrincipal) IsRoot() bool           { return false }
func (p didPrincipal) GetType() PrincipalType { return DID }

type rootPrincipal struct{}

func (p rootPrincipal) Identifier() string     { return rootIdentifier }
func (p rootPrincipal) IsAnonymous() bool      { return false }
func (p rootPrincipal) IsRoot() bool           { return true }
func (p rootPrincipal) GetType() PrincipalType { return Root }

type anonymousPrincipal struct{}

func (p anonymousPrincipal) Identifier() string     { return anonymousIdentifier }
func (p anonymousPrincipal) IsAnonymous() bool      { return true }
func (p anonymousPrincipal) IsRoot() bool           { return false }
func (p anonymousPrincipal) GetType() PrincipalType { return Anonymous }

// AnonymousPrincipal returns a generic principal representing unauthenticated users
func AnonymousPrincipal() Principal {
	return anonymousPrincipal{}
}

// RootPrincipal returns the root for the system
func RootPrincipal() Principal {
	return rootPrincipal{}
}

func NewDIDPrincipal(did string) (Principal, error) {
	err := didpkg.IsValidDID(did)
	if err != nil {
		return nil, fmt.Errorf("invalid principal id: %v", err)
	}
	return didPrincipal(did), nil
}

func ExtractPrincipal(ctx context.Context) (Principal, error) {
	any := ctx.Value(principalCtxKey)
	switch p := any.(type) {
	case nil:
		return AnonymousPrincipal(), nil
	case Principal:
		return p, nil
	default:
		return nil, ErrInvalidPrincipal
	}
}

func ExtractPrincipalWithType(ctx context.Context, t PrincipalType) (Principal, error) {
	principal, err := ExtractPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if principal.IsAnonymous() {
		return nil, ErrUnauthenticatd
	}
	if principal.GetType() != t {
		return nil, fmt.Errorf("principal type: wanted %v: got %v: %w", t, principal.GetType(), ErrPrincipalMismatch)
	}
	return principal, nil
}

func InjectPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalCtxKey, principal)
}
