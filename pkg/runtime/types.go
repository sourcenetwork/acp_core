package runtime

import (
	"context"
	"fmt"

	"cosmossdk.io/log"
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
		kv, closeFn, err := rcdb.NewLevelDB(kvPath, defualtDbFile)
		if err != nil {
			return fmt.Errorf("could not create kv: %v", err)
		}
		m.kvStore = kv
		m.cleanupFns = append(m.cleanupFns, closeFn)

		return nil
	}
}

// WithLogger sets the logger for the Runtime
func WithLogger(logger Logger) Opt {
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

// WithLogicalClock sets the LogicalClock implementation for the Runtime
func WithLogicalClock(clock LogicalClockService) Opt {
	return func(m *runtimeManager) error {
		m.clock = clock
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
		logger:   log.NewNopLogger(),
	}
	WithMemKV()(rt)

	for _, o := range opts {
		err := o(rt)
		if err != nil {
			return nil, err
		}
	}

	return rt, nil
}

type runtimeManager struct {
	kvStore    KVStore
	eventMan   EventManager
	logger     Logger
	metrics    MetricService
	clock      LogicalClockService
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

func (m *runtimeManager) GetLogicalClock() LogicalClockService {
	return m.clock
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
	GetLogicalClock() LogicalClockService

	// Terminate frees up all used up resources, leaving the runtime unusable
	Terminate() error
}

type KVStore = rcdb.KVStore

type Logger = log.Logger

type MetricService interface {
}

type EventManager interface {
	EmitEvent(any) error
}

type LogicalClockService interface {
	GetCurrentTime(ctx context.Context) uint64 // maybe make this generic over seomthing that can be compared
}
