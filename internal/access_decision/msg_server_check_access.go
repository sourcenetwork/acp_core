package access_decision

/*
import (
	"context"
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

func (k Service) CheckAccess(ctx context.Context, msg *types.MsgCheckAccess) (*types.MsgCheckAccessResponse, error) {
	/*
		runtime := k.Keeper.GetRuntime(ctx)

		repository := runtime.GetRepositoryManager().GetAccessDecisionRepository()
		paramsRepository := runtime.GetRepositoryManager().GetParamsRepository()

		engine, err := k.GetZanziEngine(ctx)
		if err != nil {
			return nil, err
		}

		record, err := engine.GetPolicy(ctx, msg.PolicyId)
		if err != nil {
			return nil, err
		}
		if record == nil {
			return nil, fmt.Errorf("policy %v: %w", msg.PolicyId, types.ErrPolicyNotFound)
		}

		creatorAddr, err := sdk.AccAddressFromBech32(msg.Creator)
		if err != nil {
			return nil, types.ErrInvalidAccAddr
		}
		creatorAcc := k.accountKeeper.GetAccount(ctx, creatorAddr)
		if creatorAcc == nil {
			return nil, types.ErrAccNotFound
		}

		cmd := access_decision.EvaluateAccessRequestsCommand{
			Policy:        record.Policy,
			Operations:    msg.AccessRequest.Operations,
			Actor:         msg.AccessRequest.Actor.Id,
			CreationTime:  msg.CreationTime,
			Creator:       creatorAcc,
			CurrentHeight: uint64(ctx.BlockHeight()),
		}
		decision, err := cmd.Execute(ctx, engine, repository, &paramsRepository)
		if err != nil {
			return nil, err
		}

		err = eventManager.EmitTypedEvent(&types.EventAccessDecisionCreated{
			Creator:    msg.Creator,
			PolicyId:   msg.PolicyId,
			DecisionId: decision.Id,
			Actor:      decision.Actor,
		})
		if err != nil {
			return nil, err
		}

		return &types.MsgCheckAccessResponse{
			Decision: decision,
		}, nil
	* /
	return nil, fmt.Errorf("not implemented")
}
*/
