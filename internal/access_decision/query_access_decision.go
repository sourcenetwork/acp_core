package access_decision

/*
import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

func (k Keeper) AccessDecision(goCtx context.Context, req *types.QueryAccessDecisionRequest) (*types.QueryAccessDecisionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	repository := k.GetAccessDecisionRepository(ctx)

	decision, err := repository.Get(goCtx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryAccessDecisionResponse{
		Decision: decision,
	}, nil
}

*/
