// package container acts a provider for service types in acp_core
package container

import (
	"github.com/sourcenetwork/acp_core/internal/params"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
)

func NewContainer(runtime *runtime.RuntimeManager) (*Container, error) {
	repo := params.NewParamsRepository(runtime)
	zanzi, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	return &Container{
		runtime:   runtime,
		paramRepo: repo,
		zanzi:     zanzi,
	}, nil
}

type Container struct {
	runtime   *runtime.RuntimeManager
	paramRepo params.ParamsRepository
	zanzi     *zanzi.Adapter
}

func (c *Container) GetParamsRepository() params.ParamsRepository {
	return c.paramRepo
}

func (c *Container) GetZanzi() *zanzi.Adapter {
	return c.zanzi
}

func (c *Container) GetLogger() runtime.Logger {
	return c.runtime.GetLogger()
}

func (c *Container) GetRequestLogger() runtime.Logger {
	return c.runtime.GetRequestLogger()
}

/*
func (c *Container) GetKVStore() KVStore {
	return

}

func (c *Container) GetEventManager() EventManager {

}


func (c *Container) GetMetricService() MetricService {

}
*/
