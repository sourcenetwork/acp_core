package did

/*
import (
	"context"

	"github.com/cyware/ssi-sdk/did"
	"github.com/cyware/ssi-sdk/did/key"
)

type DIDDocument struct {
	*did.Document
}

type Resolver interface {
	Resolve(ctx context.Context, did string) (DIDDocument, error)
}

var _ Resolver = (*KeyResolver)(nil)

type KeyResolver struct {
}

func (r *KeyResolver) Resolve(ctx context.Context, did string) (DIDDocument, error) {
	didkey := key.DIDKey(did)
	doc, err := didkey.Expand()
	if err != nil {
		return DIDDocument{}, err
	}

	return DIDDocument{
		Document: doc,
	}, nil
}

*/
