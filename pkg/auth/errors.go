package auth

import "github.com/sourcenetwork/acp_core/pkg/errors"

var ErrInvalidPrincipal = errors.New("invalid principal", errors.ErrorType_UNAUTHENTICATED)
var ErrPrincipalMismatch = errors.Wrap("principal mismatch", errors.ErrorType_UNAUTHORIZED)
