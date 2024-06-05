package did

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidDID(t *testing.T) {
	err := IsValidDID("did:key:zQ3shw2LqCxawpfYWypqSyircWhn56ZCFrzF7uNNZAnZFsA1g")
	require.Nil(t, err)
}
