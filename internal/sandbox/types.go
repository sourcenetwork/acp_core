package sandbox

import (
	"fmt"

	rcdb "github.com/sourcenetwork/raccoondb"

	"github.com/sourcenetwork/acp_core/pkg/runtime"
)

const sandboxStorePrefix = "/sandboxes"

func GetManagerForSandbox(manager runtime.RuntimeManager, handle uint64) (runtime.RuntimeManager, error) {
	kv := manager.GetKVStore()
	prefix := fmt.Sprintf("/engine/%v", handle)
	kv = rcdb.NewWrapperKV(kv, []byte(prefix))
	engineManager, err := runtime.NewRuntimeManager(
		runtime.WithKVStore(kv),
		runtime.WithLogger(manager.GetLogger()),
	)
	if err != nil {
		return nil, err
	}
	return engineManager, nil
}
