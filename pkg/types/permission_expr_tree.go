package types

import "github.com/cosmos/gogoproto/jsonpb"

func (t *PermissionFetchTree) MarshalJSON() (string, error) {
	marshaler := jsonpb.Marshaler{}
	return marshaler.MarshalToString(t)
}

func (t *PermissionFetchTree) UnmarshalJSON(json string) error {
	return jsonpb.UnmarshalString(json, t)
}
