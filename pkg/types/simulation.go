package types

// HasErrors return true if DeclarationErrors contains any errors
func (e *DeclarationErrors) HasErrors() bool {
	return e.PolicyErrors != nil && len(e.PolicyErrors) > 0 && e.RelationshipsErrors != nil && len(e.RelationshipsErrors) > 0 && e.TheoremsErrrors != nil && len(e.TheoremsErrrors) > 0
}
