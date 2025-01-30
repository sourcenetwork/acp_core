package types

func NewOperation(obj *Object, permission string) *Operation {
	return &Operation{
		Object:     obj,
		Permission: permission,
	}
}
