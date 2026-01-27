package runtime_test

import (
	"context"
	"errors"
	"pan/internal/runtime"
	"reflect"
	"testing"

	mocked "pan/mocks/pan/internal_/runtime"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestModule struct {
}

func TestEngine(t *testing.T) {

	getRegistryImplTypeName := func() string {
		rtype := reflect.TypeOf(runtime.NewRegistry())
		return rtype.String()
	}

	setup := func() (engine *runtime.Engine) {
		engine, _ = runtime.New()
		return
	}

	t.Run("Mount", func(t *testing.T) {
		e := setup()

		module := &mocked.MockModule{}
		defer module.AssertExpectations(t)
		module.On("TypeOfModule").Once().Return([]reflect.Type{reflect.TypeOf(module)})

		err := e.Mount(module, &TestModule{})
		assert.Nil(t, err)

	})

	t.Run("Mount with provider module", func(t *testing.T) {
		e := setup()

		module := &mocked.MockModule{}
		defer module.AssertExpectations(t)
		module.On("TypeOfModule").Once().Return([]reflect.Type{reflect.TypeOf(module)})

		providerModule := &mocked.MockProviderModule{}
		defer providerModule.AssertExpectations(t)
		providerModule.On("Modules").Once().Return([]interface{}{module, &TestModule{}})

		err := e.Mount(providerModule)
		assert.Nil(t, err)
	})

	t.Run("Mount with initialize module", func(t *testing.T) {
		e := setup()
		registryTypeName := getRegistryImplTypeName()
		ctx := context.Background()

		initializeModule := &mocked.MockInitializeModule{}
		defer initializeModule.AssertExpectations(t)
		initErr := errors.New("test error")
		initializeModule.On("Init", ctx, mock.AnythingOfType(registryTypeName)).Once().Return(initErr)

		err := e.Mount(initializeModule)
		assert.Nil(t, err)

		err = e.Bootstrap(ctx)
		assert.Equal(t, initErr, err)

		initializeModule.On("Init", ctx, mock.AnythingOfType(registryTypeName)).Once().Return(nil)

		err = e.Bootstrap(context.Background())
		assert.Nil(t, err)
	})

}
