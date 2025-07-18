package zanzi

import (
	"context"
	"fmt"

	rcdb "github.com/sourcenetwork/raccoondb"
	"github.com/sourcenetwork/zanzi"
	"github.com/sourcenetwork/zanzi/pkg/api"
	"github.com/sourcenetwork/zanzi/pkg/domain"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

// NewZanzi builds an AuthEngine with zanzi as backend
func NewZanzi(kv rcdb.KVStore, logger types.Logger) (*Adapter, error) {
	wrappedLogger := &loggerWrapper{}

	z, err := zanzi.New(
		zanzi.WithKVStore(kv),
		zanzi.WithLogger(wrappedLogger),
	)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		zanzi:        z,
		policyMapper: policyMapper{},
	}, nil
}

// Adapter wraps Zanzi's API Service and adapts acp_cores domain model to zanzi's
type Adapter struct {
	zanzi        zanzi.Zanzi
	policyMapper policyMapper
}

// Reurn a Relationship from a Policy, returns nil if Relationship does not exist
func (z *Adapter) GetRelationship(ctx context.Context, policy *types.Policy, rel *types.Relationship) (*types.RelationshipRecord, error) {
	serv := z.zanzi.GetPolicyService()
	mapper := newRelationshipMapper(policy.ActorResource.Name)

	req := &api.GetRelationshipRequest{
		PolicyId:     policy.Id,
		Relationship: mapper.ToZanziRelationship(rel),
	}

	result, err := serv.GetRelationship(ctx, req)
	err = mapErr(err)
	if err != nil {
		return nil, fmt.Errorf("GetRelationship: %w", err)
	}

	fetchedRel, err := mapper.FromZanziRelationship(result.Record)
	if err != nil {
		return nil, fmt.Errorf("GetRelationship: %w", err)
	}

	return fetchedRel, nil
}

func (z *Adapter) ValidateRelationship(ctx context.Context, policy *types.Policy, rel *types.Relationship) (valid bool, msg string, err error) {
	serv := z.zanzi.GetPolicyService()
	mapper := newRelationshipMapper(policy.ActorResource.Name)

	req := &api.ValidateRelationshipRequest{
		PolicyId:     policy.Id,
		Relationship: mapper.ToZanziRelationship(rel),
	}

	result, err := serv.ValidateRelationship(ctx, req)
	err = mapErr(err)
	if err != nil {
		return false, "", fmt.Errorf("ValidateRelationship: %w", err)
	}

	return result.Valid, result.ErrorMsg, nil
}

// Sets a Relationship within a Policy
func (z *Adapter) SetRelationship(ctx context.Context, policy *types.Policy, rec *types.RelationshipRecord) (RecordFound, error) {
	serv := z.zanzi.GetPolicyService()
	mapper := newRelationshipMapper(policy.ActorResource.Name)

	rec.PolicyId = policy.Id
	zanziRecord, err := mapper.ToZanziRelationshipRecord(rec)
	if err != nil {
		return false, fmt.Errorf("SetRelationship: %w", err)
	}

	req := &api.SetRelationshipRequest{
		PolicyId:     policy.Id,
		Relationship: zanziRecord.Relationship,
		AppData:      zanziRecord.AppData,
	}

	response, err := serv.SetRelationship(ctx, req)
	err = mapErr(err)
	if err != nil {
		return false, fmt.Errorf("SetRelationship: %w", err)
	}

	return RecordFound(response.RecordOverwritten), nil
}

// GetPolicy returns a PolicyRecord for the given id
func (z *Adapter) GetPolicy(ctx context.Context, policyId string) (*types.PolicyRecord, error) {
	serv := z.zanzi.GetPolicyService()

	req := api.GetPolicyRequest{
		Id: policyId,
	}
	res, err := serv.GetPolicy(ctx, &req)
	err = mapErr(err)
	if err != nil {
		return nil, err
	}
	if res.Record == nil {
		return nil, nil
	}

	mapped, err := z.policyMapper.FromZanzi(res.Record)
	if err != nil {
		return nil, err
	}

	return mapped, nil
}

// SetPolicy stores a new Policy with the given Id
func (z *Adapter) SetPolicy(ctx context.Context, record *types.PolicyRecord) error {
	serv := z.zanzi.GetPolicyService()

	zanziRecord, err := z.policyMapper.ToZanziRecord(record)
	if err != nil {
		return err
	}

	req := api.CreatePolicyRequest{
		PolicyDefinition: &api.PolicyDefinition{
			Definition: &api.PolicyDefinition_Policy{
				Policy: zanziRecord.Policy,
			},
		},
		AppData: zanziRecord.AppData,
	}
	_, err = serv.CreatePolicy(ctx, &req)
	err = mapErr(err)
	if err != nil {
		return err
	}

	return nil
}

func (z *Adapter) EditPolicy(ctx context.Context, record *types.PolicyRecord) (uint64, error) {
	serv := z.zanzi.GetPolicyService()

	zanziRecord, err := z.policyMapper.ToZanziRecord(record)
	if err != nil {
		return 0, err
	}

	req := api.EditPolicyRequest{
		PolicyId: record.Policy.Id,
		PolicyDefinition: &api.PolicyDefinition{
			Definition: &api.PolicyDefinition_Policy{
				Policy: zanziRecord.Policy,
			},
		},
	}
	result, err := serv.EditPolicy(ctx, &req)
	err = mapErr(err)
	if err != nil {
		return 0, err
	}

	_, err = serv.EditPolicyAppData(ctx, &api.EditPolicyAppDataRequest{
		PolicyId: record.Policy.Id,
		AppData:  zanziRecord.AppData,
	})
	err = mapErr(err)
	if err != nil {
		return 0, err
	}

	return result.RemovedRelationshipsCount, nil
}

// Returns all Relationships which matches selector
func (z *Adapter) FilterRelationships(ctx context.Context, policy *types.Policy, selector *types.RelationshipSelector) ([]*types.RelationshipRecord, error) {
	serv := z.zanzi.GetPolicyService()
	relationshipMapper := newRelationshipMapper(policy.ActorResource.Name)
	selectorMapper := newSelectorMapper(relationshipMapper)

	zanziSelector, err := selectorMapper.ToZanziSelector(selector)
	if err != nil {
		return nil, fmt.Errorf("FilterRelationships: %v", err)
	}

	req := api.FindRelationshipRecordsRequest{
		PolicyId: policy.Id,
		Selector: zanziSelector,
	}

	resp, err := serv.FindRelationshipRecords(ctx, &req)
	err = mapErr(err)
	if err != nil {
		return nil, fmt.Errorf("FilterRelationships: %v", err)
	}

	records, err := utils.MapFailableSlice(resp.Result.Records, relationshipMapper.FromZanziRelationship)
	if err != nil {
		return nil, fmt.Errorf("FilterRelationships: %v", err)
	}

	return records, nil
}

// Check verifies whether an Acccess Request is allowed within a certain Policy
func (z *Adapter) Check(ctx context.Context, policy *types.Policy, operation *types.Operation, actor *types.Actor) (bool, error) {
	service := z.zanzi.GetRelationGraphService()
	mapper := newRelationshipMapper(policy.ActorResource.Name)

	req := &api.CheckRequest{
		PolicyId: policy.Id,
		AccessRequest: &domain.AccessRequest{
			Object:   mapper.MapObject(operation.Object),
			Relation: operation.Permission,
			Subject: &domain.Entity{
				Resource: policy.ActorResource.Name,
				Id:       actor.Id,
			},
		},
	}
	response, err := service.Check(ctx, req)
	err = mapErr(err)
	if err != nil {
		return false, fmt.Errorf("Check: %w", err)
	}

	return response.Result.Authorized, nil
}

// DeleteRelationship removes a Relationship from a Policy
func (z *Adapter) DeleteRelationship(ctx context.Context, policy *types.Policy, relationship *types.Relationship) (RecordFound, error) {
	service := z.zanzi.GetPolicyService()
	mapper := newRelationshipMapper(policy.ActorResource.Name)

	req := api.DeleteRelationshipRequest{
		PolicyId:     policy.Id,
		Relationship: mapper.ToZanziRelationship(relationship),
	}
	response, err := service.DeleteRelationship(ctx, &req)
	err = mapErr(err)
	if err != nil {
		return false, err
	}

	return RecordFound(response.Found), nil
}

// DeleteRelationships removes all Relationships matching the given selector
func (z *Adapter) DeleteRelationships(ctx context.Context, policy *types.Policy, selector *types.RelationshipSelector) (uint, error) {
	service := z.zanzi.GetPolicyService()
	relationshipMapper := newRelationshipMapper(policy.ActorResource.Name)
	selectorMapper := newSelectorMapper(relationshipMapper)

	zanziSelector, err := selectorMapper.ToZanziSelector(selector)
	if err != nil {
		return 0, fmt.Errorf("DeleteRelationships: %w", err)
	}

	request := api.DeleteRelationshipsRequest{
		PolicyId: policy.Id,
		Selector: zanziSelector,
	}
	response, err := service.DeleteRelationships(ctx, &request)
	err = mapErr(err)
	if err != nil {
		return 0, fmt.Errorf("DeleteRelationships: %w", err)
	}

	return uint(response.RecordsAffected), nil
}

// ListPolicyIds returns the IDs of all known Policies
func (z *Adapter) ListPolicyIds(ctx context.Context) ([]string, error) {
	service := z.zanzi.GetPolicyService()

	req := api.ListPolicyIdsRequest{}
	resp, err := service.ListPolicyIds(ctx, &req)
	err = mapErr(err)
	if err != nil {
		return nil, fmt.Errorf("ListPolicyIds: %w", err)
	}

	// Use MapNullableSlice instead of MapSlice to filter out 'nil' policy ids.
	return utils.MapNullableSlice(resp.Records, func(rec *api.ListPolicyIdsResponse_Record) string {
		return rec.Id
	}), nil
}

func (z *Adapter) ListPolicies(ctx context.Context) ([]*types.PolicyRecord, error) {
	resp, err := z.zanzi.GetPolicyService().ListPolicies(ctx, &api.ListPoliciesRequest{})
	err = mapErr(err)
	if err != nil {
		return nil, fmt.Errorf("ListPolicies: %v", err)
	}

	// This filter fixes a bug which, strangely, was only observable from SourceHub.
	// Sometimes zanzi would return a PolicyRecord with Metadat but whose rec.Policy == nil.
	// The root cause is unknown but for now that fixes the issue
	zanziRecords := utils.FilterSlice(resp.Records, func(rec *domain.PolicyRecord) bool { return rec != nil && rec.Policy != nil })

	records, err := utils.MapFailableSlice(zanziRecords, func(rec *domain.PolicyRecord) (*types.PolicyRecord, error) {
		return z.policyMapper.FromZanzi(rec)
	})
	if err != nil {
		return nil, fmt.Errorf("ListPolicies: %v", err)
	}

	return records, nil
}

func (z *Adapter) ValidatePolicy(ctx context.Context, policy *types.Policy) (valid bool, msg string, err error) {
	zanziPolicy := z.policyMapper.ToZanzi(policy)

	req := api.ValidatePolicyRequest{
		PolicyDefinition: &api.PolicyDefinition{
			Definition: &api.PolicyDefinition_Policy{
				Policy: zanziPolicy,
			},
		},
	}
	response, err := z.zanzi.GetPolicyService().ValdiatePolicy(ctx, &req)
	err = mapErr(err)
	valid = response.Valid
	msg = response.ErrorMsg

	return
}

func (z *Adapter) DeletePolicy(ctx context.Context, id string) (bool, error) {
	resp, err := z.zanzi.GetPolicyService().DeletePolicy(ctx, &api.DeletePolicyRequest{Id: id})
	err = mapErr(err)
	if err != nil {
		return false, err
	}

	return resp.Found, nil
}

func (z *Adapter) Expand(ctx context.Context, policy *types.Policy, object *types.Object, relation string) (string, error) {
	service := z.zanzi.GetRelationGraphService()

	req := &api.ExpandRequest{
		PolicyId: policy.Id,
		Root: &domain.RelationNode{
			Node: &domain.RelationNode_EntitySet{
				EntitySet: &domain.EntitySetNode{
					Object: &domain.Entity{
						Resource: object.Resource,
						Id:       object.Id,
					},
					Relation: relation,
				},
			},
		},
		Format: api.ExplainFormat_DOT,
	}
	response, err := service.Expand(ctx, req)
	err = mapErr(err)
	if err != nil {
		return "", fmt.Errorf("Expand: %w", err)
	}

	return response.GoalTree, nil
}

// ExplainCheck verifies whether an Acccess Request is allowed within a certain Policy
func (z *Adapter) ExplainCheck(ctx context.Context, policy *types.Policy, operation *types.Operation, actor *types.Actor) (bool, string, error) {
	service := z.zanzi.GetRelationGraphService()
	mapper := newRelationshipMapper(policy.ActorResource.Name)

	req := &api.ExplainCheckRequest{
		PolicyId: policy.Id,
		AccessRequest: &domain.AccessRequest{
			Object:   mapper.MapObject(operation.Object),
			Relation: operation.Permission,
			Subject: &domain.Entity{
				Resource: policy.ActorResource.Name,
				Id:       actor.Id,
			},
		},
		Format: api.ExplainFormat_DOT,
	}
	response, err := service.ExplainCheck(ctx, req)
	err = mapErr(err)
	if err != nil {
		return false, "", fmt.Errorf("ExplainCheck: %w", err)
	}

	return response.Authorized, response.GoalTree, nil
}
