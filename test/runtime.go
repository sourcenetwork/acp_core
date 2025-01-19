package test

import (
	"context"
	"testing"

	prototypes "github.com/cosmos/gogoproto/types"

	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/stretchr/testify/require"
)

// NewTestRuntime returns a runtime for executing tests
func NewTestRuntime(t testing.TB, timeServ runtime.TimeService) runtime.RuntimeManager {
	manager, err := runtime.NewRuntimeManager(
		runtime.WithMemKV(),
		runtime.WithTimeService(timeServ),
	)
	require.Nil(t, err)
	return manager
}

var _ runtime.TimeService = (*constantTimeService)(nil)

func NewConstantTimeService(ts *prototypes.Timestamp) runtime.TimeService {
	return &constantTimeService{
		ts: ts,
	}
}

type constantTimeService struct {
	ts *prototypes.Timestamp
}

func (s *constantTimeService) GetNow(_ context.Context) (*prototypes.Timestamp, error) {
	return s.ts, nil
}
