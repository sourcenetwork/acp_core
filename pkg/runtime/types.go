package runtime

import (
	"context"
	"fmt"
	"time"

	prototypes "github.com/cosmos/gogoproto/types"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	rcdb "github.com/sourcenetwork/raccoondb"
)

var _ RuntimeManager = (*runtimeManager)(nil)

type AccountDIDIssuer interface {
	IssueDID(ctx context.Context, id string) (string, error)
}

type Opt func(m *runtimeManager) error

// WithKVStore specifies the kKVStore which must be used
func WithKVStore(kv KVStore) Opt {
	return func(m *runtimeManager) error {
		m.kvStore = kv
		return nil
	}
}

// WithPersistentKV initiates a LevelDB store at the given path
// and uses it as the default KV Store
func WithPersistentKV(kvPath string) Opt {
	return func(m *runtimeManager) error {
		kv, closeFn, err := rcdb.NewPersistentKV(kvPath, defualtDbFile)
		if err != nil {
			return fmt.Errorf("could not create kv: %v", err)
		}
		m.kvStore = kv
		m.cleanupFns = append(m.cleanupFns, closeFn)

		return nil
	}
}

// WithLogger sets the logger for the Runtime
func WithLogger(logger types.Logger) Opt {
	return func(m *runtimeManager) error {
		m.logger = logger
		return nil
	}
}

// WithEventManager sets the EventManager for the Runtime
func WithEventManager(manager EventManager) Opt {
	return func(m *runtimeManager) error {
		m.eventMan = manager
		return nil
	}
}

// WithTimeService sets the TimeService implementation for the Runtime
func WithTimeService(service TimeService) Opt {
	return func(m *runtimeManager) error {
		m.timeServ = service
		return nil
	}
}

// WithMemKV initiates an memory kv store and sets it as the Runtime's KV
func WithMemKV() Opt {
	return func(m *runtimeManager) error {
		m.kvStore = rcdb.NewMemKV()
		return nil
	}
}

// NewRuntimeManager creates a RuntimeManager with the given options
func NewRuntimeManager(opts ...Opt) (RuntimeManager, error) {
	rt := &runtimeManager{
		eventMan: &DefaultEventManager{},
		logger:   types.NoopLogger(),
		memKV:    rcdb.NewMemKV(),
		timeServ: NewLocalTimeService(),
	}
	WithMemKV()(rt)

	for _, o := range opts {
		err := o(rt)
		if err != nil {
			return nil, errors.NewFromBaseError(err, errors.ErrorType_INTERNAL, "building runtime manager")
		}
	}

	return rt, nil
}

type runtimeManager struct {
	kvStore    KVStore
	memKV      KVStore
	eventMan   EventManager
	logger     Logger
	metrics    MetricService
	timeServ   TimeService
	cleanupFns []func() error
	terminated bool
}

func (m *runtimeManager) GetKVStore() KVStore {
	return m.kvStore
}

func (m *runtimeManager) GetEventManager() EventManager {
	return m.eventMan
}

func (m *runtimeManager) GetLogger() Logger {
	return m.logger
}

func (m *runtimeManager) GetMetricService() MetricService {
	return m.metrics
}

func (m *runtimeManager) GetAccountDIDIssuer() AccountDIDIssuer {
	return nil
}

func (m *runtimeManager) GetTimeService() TimeService {
	return m.timeServ
}

func (m *runtimeManager) Terminate() error {
	m.terminated = true
	for _, f := range m.cleanupFns {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}

type RuntimeManager interface {
	GetKVStore() KVStore
	GetEventManager() EventManager
	GetLogger() Logger
	GetMetricService() MetricService
	GetTimeService() TimeService

	// Terminate frees up all used up resources, leaving the runtime unusable
	Terminate() error
}

type KVStore = rcdb.KVStore

type Logger = types.Logger

type MetricService interface {
}

type EventManager interface {
	EmitEvent(any) error
}

// TimeService models a Time oracle
type TimeService interface {
	// GetNow returns the current Time
	GetNow(ctx context.Context) (*prototypes.Timestamp, error)
}

var _ TimeService = (*LocalTimeService)(nil)

// NewLocalTimeService returns a TimeService implementation which uses the local clock to return the time
// This is the default TimeService implementation used in the RuntimeManager
func NewLocalTimeService() TimeService {
	return &LocalTimeService{}
}

// LocalTimeService returns time from the local clock
type LocalTimeService struct{}

func (s *LocalTimeService) GetNow(ctx context.Context) (*prototypes.Timestamp, error) {
	now := time.Now()
	ts, err := prototypes.TimestampProto(now)
	if err != nil {
		return nil, errors.NewFromBaseError(err, errors.ErrorType_INTERNAL, "LocalTimeService failed: converting timestamp")
	}
	return ts, nil
}
