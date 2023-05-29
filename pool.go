package mjml

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/puddle/v2"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

func constructor(ctx context.Context) (api.Module, error) {

	id, err := randomIdentifier()

	if err != nil {
		return nil, fmt.Errorf("error generating random id for wasm module: %w", err)
	}

	idStr := strconv.Itoa(int(id))

	module, err := runtime.InstantiateModule(ctx, compiled, wazero.NewModuleConfig().WithName(idStr))

	if err != nil {
		return nil, fmt.Errorf("error instantiating wasm module: %w", err)
	}

	return module, nil
}

func destructor(module api.Module) {
	_ = module.Close(context.Background()) // Not possible to deal with this error
}

func newResourcePool(maxSize int32) (*puddle.Pool[api.Module], error) {
	//pool := puddle.NewPool(constructor, destructor, maxSize)
	pool, err := puddle.NewPool(&puddle.Config[api.Module]{Constructor: constructor, Destructor: destructor, MaxSize: maxSize})

	if err != nil {
		return pool, fmt.Errorf("error creating resource pool: %w", err)
	}

	err = pool.CreateResource(context.Background())

	if err != nil {
		return pool, fmt.Errorf("error prewarming resource pool: %w", err)
	}

	return pool, nil
}

func periodicallyRemoveIdleResources(pool *puddle.Pool[api.Module]) {

	duration := 2 * time.Second
	ticker := time.NewTicker(duration)

	for {
		select {
		case <-ticker.C:
			stats := pool.Stat()

			if stats.TotalResources() <= 1 {
				continue
			}

			idleResources := pool.AcquireAllIdle()
			numIdleResources := len(idleResources)

			if numIdleResources <= 0 {
				continue
			}

			max := int(stats.TotalResources())

			amountToKill := numIdleResources

			if numIdleResources >= max {
				amountToKill = numIdleResources - 1
			}

			for i := 0; i < numIdleResources; i++ {
				if i >= amountToKill {
					idleResources[i].Release()
				} else {
					idleResources[i].Destroy()
				}
			}
		}
	}
}
