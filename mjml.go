package mjml

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/jackc/puddle"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed wasm/mjml.wasm.br
var wasm []byte

var (
	runtime      wazero.Runtime
	compiled     wazero.CompiledModule
	results      *sync.Map
	resourcePool *puddle.Pool
)

func init() {
	ctx := context.Background()

	results = &sync.Map{}

	br := brotli.NewReader(bytes.NewReader(wasm))
	decompressed, err := io.ReadAll(br)

	if err != nil {
		panic(fmt.Sprintf("Error decompressing wasm file: %s", err))
	}

	runtime = wazero.NewRuntime(ctx) // TODO: this should be closed

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, runtime); err != nil {
		panic(fmt.Sprintf("Error instantiating wasi snapshot preview 1: %s", err))
	}

	err = registerHostFunctions(ctx, runtime)

	if err != nil {
		panic(fmt.Sprintf("Error registering host functions: %s", err))
	}

	compiled, err = runtime.CompileModule(ctx, decompressed, wazero.NewCompileConfig())

	if err != nil {
		panic(fmt.Sprintf("Error compiling wasm module: %s", err))
	}

	resourcePool, err = newResourcePool(10)

	if err != nil {
		panic(fmt.Sprintf("Error creating resource pool: %s", err))
	}

	go periodicallyRemoveIdleResources(resourcePool)
}

func SetMaxWorkers(maxSize int32) error {
	oldPool := resourcePool

	newPool, err := newResourcePool(maxSize)

	if err != nil {
		return fmt.Errorf("error creating new resource pool: %w", err)
	}

	resourcePool = newPool
	oldPool.Close()

	return nil
}

type jsonResult struct {
	HTML  string `json:"html"`
	Error *Error `json:"error,omitempty"`
}

// ToHTML converts a string containing mjml to HTML while using any of the optionally provided options
func ToHTML(ctx context.Context, mjml string, toHTMLOptions ...ToHTMLOption) (string, error) {
	data := map[string]interface{}{
		"mjml": mjml,
	}

	o := options{
		data: map[string]interface{}{},
	}

	for _, opt := range toHTMLOptions {
		opt(o)
	}

	if len(o.data) > 0 {
		data["options"] = o.data
	}

	inputBytes := bytes.NewBuffer([]byte{})

	encoder := json.NewEncoder(inputBytes)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(data)

	if err != nil {
		return "", fmt.Errorf("error encoding input data: %w", err)
	}

	jsonInput := inputBytes.String()
	jsonInputLen := uint64(len(jsonInput))

	var (
		module *puddle.Resource
		tries  int
	)

	for {
		tries++

		var err error

		module, err = resourcePool.Acquire(ctx)

		if err != nil {

			if tries >= 30 {
				return "", fmt.Errorf("unable to accquire wasm module after 30 tries: %w", err)
			}

			if err == puddle.ErrClosedPool {
				time.Sleep(1 * time.Millisecond)
				continue
			}

			return "", fmt.Errorf("error accquiring wasm module: %w", err)
		}

		break
	}

	defer module.Release()

	mod, ok := module.Value().(api.Module)

	if !ok {
		return "", errors.New("pool resource is not an api.Module")
	}

	deallocate := mod.ExportedFunction("deallocate")
	allocate := mod.ExportedFunction("allocate")
	run := mod.ExportedFunction("run_e")
	memory := mod.Memory()

	allocation, err := allocate.Call(ctx, jsonInputLen)

	if err != nil {
		return "", fmt.Errorf("error allocating memory: %w", err)
	}

	if len(allocation) != 1 {
		return "", errors.New("invalid input pointer allocated")
	}

	inputPtr := allocation[0]

	defer deallocate.Call(ctx, inputPtr)

	if !memory.Write(ctx, uint32(inputPtr), []byte(jsonInput)) {
		return "", fmt.Errorf("error writing input to memory: %w", err)
	}

	ident, err := randomIdentifier()

	if err != nil {
		return "", fmt.Errorf("error generating identifier: %w", err)
	}

	resultCh := make(chan []byte, 1)

	results.Store(ident, resultCh)

	defer results.Delete(ident)

	_, err = run.Call(ctx, inputPtr, jsonInputLen, uint64(ident))

	if err != nil {
		return "", fmt.Errorf("error calling run: %w", err)
	}

	result := <-resultCh

	res := jsonResult{}

	err = json.Unmarshal(result, &res)

	if err != nil {
		return "", fmt.Errorf("error decoding result json: %w", err)
	}

	if res.Error != nil {
		return "", *res.Error
	}

	return res.HTML, nil
}

func registerHostFunctions(ctx context.Context, r wazero.Runtime) error {

	_, err := r.NewModuleBuilder("env").
		ExportFunction("return_result", func(ctx context.Context, m api.Module, ptr uint32, len uint32, ident uint32) {
			if ch, ok := results.Load(int32(ident)); ok {

				result, ok := m.Memory().Read(ctx, ptr, len)

				resultCh, isResultCh := ch.(chan []byte)

				if ok && isResultCh {
					resultCh <- result
				}
			}
		}).
		ExportFunction("get_static_file", func(_ uint32, _ uint32, _ uint32) uint32 {
			panic("get_static_file is unimplemented")
		}).
		ExportFunction("request_set_field", func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("request_set_field is unimplemented")
		}).
		ExportFunction("resp_set_header", func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32) {
			panic("resp_set_header is unimplemented")
		}).
		ExportFunction("cache_get", func(_ uint32, _ uint32, _ uint32) uint32 {
			panic("cache_get is unimplemented")
		}).
		ExportFunction("add_ffi_var", func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("add_ffi_var is unimplemented")
		}).
		ExportFunction("get_ffi_result", func(_ uint32, _ uint32) uint32 {
			panic("get_ffi_result is unimplemented")
		}).
		ExportFunction("return_error", func(_ uint32, _ uint32, _ uint32, _ uint32) {
			panic("return_error is unimplemented")
		}).
		ExportFunction("fetch_url", func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("fetch_url is unimplemented")
		}).
		ExportFunction("graphql_query", func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("graphql_query is unimplemented")
		}).
		ExportFunction("db_exec", func(_ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("db_exec is unimplemented")
		}).
		ExportFunction("cache_set", func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("cache_set is unimplemented")
		}).
		ExportFunction("request_get_field", func(_ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("request_get_field is unimplemented")
		}).
		ExportFunction("log_msg", func(ctx context.Context, m api.Module, ptr uint32, size uint32, level uint32, ident uint32) {
			panic("log_msg is unimplemented")
		}).Instantiate(ctx, r)

	return err
}
