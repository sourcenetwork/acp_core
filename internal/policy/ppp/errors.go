package ppp

import "github.com/sourcenetwork/acp_core/pkg/errors"

var (
	ErrDiscretionaryTransformer  = errors.New("discretionary policy transformer", errors.ErrorType_BAD_INPUT)
	ErrIdTransformer             = errors.New("id policy transformer", errors.ErrorType_BAD_INPUT)
	ErrBasicTransformer          = errors.New("basic transformer", errors.ErrorType_BAD_INPUT)
	ErrAdministrationTransformer = errors.New("decentralized administration transformer", errors.ErrorType_BAD_INPUT)
	ErrDefraSpec                 = errors.New("defra policy specification", errors.ErrorType_BAD_INPUT)
	ErrPolicyProcessing          = errors.New("policy processing", errors.ErrorType_BAD_INPUT)
)
