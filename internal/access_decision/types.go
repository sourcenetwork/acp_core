package access_decision

import (
	"context"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

type Repository interface {
	Set(ctx context.Context, decision *types.AccessDecision) error

	Get(ctx context.Context, id string) (*types.AccessDecision, error)

	Delete(ctx context.Context, id string) error

	// List of Ids of all Decisions
	ListIds(ctx context.Context) ([]string, error)

	List(ctx context.Context) ([]*types.AccessDecision, error)
}

type ParamsRepository interface {
	GetDefaults(context.Context) (*types.DecisionParams, error)
}
