package runtime_test

import (
	"errors"
	"pan/runtime"
	"reflect"
	"testing"

	mocked "pan/mocks/pan/runtime"

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
		engine = runtime.New()
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
		providerModule.On("GetModules").Once().Return([]interface{}{module, &TestModule{}})

		err := e.Mount(providerModule)
		assert.Nil(t, err)
	})

	t.Run("Mount with initialize module", func(t *testing.T) {
		e := setup()
		registryTypeName := getRegistryImplTypeName()

		initializeModule := &mocked.MockInitializeModule{}
		defer initializeModule.AssertExpectations(t)
		initErr := errors.New("test error")
		initializeModule.On("Init", mock.AnythingOfType(registryTypeName)).Once().Return(initErr)

		err := e.Mount(initializeModule)
		assert.Nil(t, err)
		_, err = e.Bootstrap()
		assert.Equal(t, initErr, err)

		initializeModule.On("Init", mock.AnythingOfType(registryTypeName)).Once().Return(nil)

		_, err = e.Bootstrap()
		assert.Nil(t, err)
	})

	t.Run("Mount with ready module", func(t *testing.T) {
		e := setup()

		readyModule := &mocked.MockReadyModule{}
		defer readyModule.AssertExpectations(t)
		readyErr := errors.New("test error")
		readyModule.On("Ready").Once().Return(readyErr)

		err := e.Mount(readyModule)
		assert.Nil(t, err)
		ctx, err := e.Bootstrap()
		assert.Nil(t, err)
		err = ctx.Wait()
		assert.Equal(t, readyErr, err)

		readyModule.On("Ready").Once().Return(nil)

		ctx, err = e.Bootstrap()
		assert.Nil(t, err)
		err = ctx.Wait()
		assert.Nil(t, err)
	})
}
