package access_decision

/*
import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	//"github.com/sourcenetwork/acp_core/access_decision"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func (k Keeper) VerifyAccessRequest(goCtx context.Context, req *types.QueryVerifyAccessRequestRequest) (*types.QueryVerifyAccessRequestResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	engine, err := k.GetZanziEngine(ctx)
	if err != nil {
		return nil, err
	}

	_, err = engine.GetPolicy(goCtx, req.PolicyId)
	if err != nil {
		return nil, err
	}

	/*
		cmd := access_decision.VerifyAccessRequestQuery{
			Policy:        rec.Policy,
			AccessRequest: req.AccessRequest,
		}
		valid, err := cmd.Execute(ctx, engine)
		if err != nil {
			return nil, err
		}
	* /

	return &types.QueryVerifyAccessRequestResponse{
		Valid: true,
		//Valid: valid,
	}, nil
}

*/
