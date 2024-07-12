package runtime

import (
	"fmt"

	rcdb "github.com/sourcenetwork/raccoondb"
)

type Opt func(m *RuntimeManager) error

// WithKVStore specifies the kKVStore which must be used
func WithKVStore(kv KVStore) Opt {
	return func(m *RuntimeManager) error {
		m.kvStore = kv
		return nil
	}
}

// WithPersistentKV initiates a LevelDB store at the given path
// and uses it as the default KV Store
func WithPersistentKV(kvPath string) Opt {
	return func(m *RuntimeManager) error {
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
func WithLogger(logger Logger) Opt {
	return func(m *RuntimeManager) error {
		m.logger = logger
		return nil
	}
}

// WithEventManager sets the EventManager for the Runtime
func WithEventManager(manager EventManager) Opt {
	return func(m *RuntimeManager) error {
		m.eventMan = manager
		return nil
	}
}

// WithMemKV initiates an memory kv store and sets it as the Runtime's KV
func WithMemKV() Opt {
	return func(m *RuntimeManager) error {
		m.kvStore = rcdb.NewMemKV()
		return nil
	}
}

// NewRuntimeManager creates a RuntimeManager with the given options
func NewRuntimeManager(opts ...Opt) (*RuntimeManager, error) {
	rt := &RuntimeManager{
		eventMan: &DefaultEventManager{},
		logger:   nil,
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

type RuntimeManager struct {
	kvStore    KVStore
	eventMan   EventManager
	logger     Logger
	metrics    MetricService
	cleanupFns []func() error
	terminated bool
}

func (m *RuntimeManager) GetKVStore() KVStore {
	return m.kvStore
}

func (m *RuntimeManager) GetEventManager() EventManager {
	return m.eventMan
}

func (m *RuntimeManager) GetRequestLogger() Logger {
	return m.logger
}

func (m *RuntimeManager) GetLogger() Logger {
	return m.logger
}

func (m *RuntimeManager) GetMetricService() MetricService {
	return m.metrics
}

func (m *RuntimeManager) Terminate() error {
	m.terminated = true
	for _, f := range m.cleanupFns {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}

type KVStore = rcdb.KVStore

type MetricService interface {
}

type EventManager interface {
	EmitEvent(any) error
}
