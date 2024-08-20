package types

// HasErrors return true if DeclarationErrors contains any errors
func (e *SandboxDataErrors) HasErrors() bool {
	if e == nil {
		return false
	}

	return !((e.PolicyErrors == nil || len(e.PolicyErrors) == 0) && (e.RelationshipsErrors == nil || len(e.RelationshipsErrors) == 0) && (e.TheoremsErrors == nil || len(e.TheoremsErrors) == 0))
}
