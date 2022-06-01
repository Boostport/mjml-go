package mjml

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/puddle"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

func constructor(ctx context.Context) (interface{}, error) {

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

func destructor(module interface{}) {

	mod, ok := module.(api.Module)

	if !ok {
		panic("destructor was not given an api.Module")
	}

	_ = mod.Close(context.Background()) // Not possible to deal with this error
}

func newResourcePool(maxSize int32) (*puddle.Pool, error) {
	pool := puddle.NewPool(constructor, destructor, maxSize)

	err := pool.CreateResource(context.Background())

	if err != nil {
		return pool, fmt.Errorf("error prewarming resource pool: %w", err)
	}

	return pool, nil
}

func periodicallyRemoveIdleResources(pool *puddle.Pool) {

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
