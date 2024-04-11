package runtime_test

import (
	"pan/runtime"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {

	setup := func() runtime.Registry {
		registry := runtime.NewRegistry()
		return registry
	}

	t.Run("Append", func(t *testing.T) {
		registry := setup()

		module := &TestModule{}
		err := registry.Append(module, reflect.TypeOf(module))
		assert.Nil(t, err)

		err = registry.Append(module, reflect.TypeFor[runtime.Module]())
		assert.Equal(t, runtime.ErrModuleType, err)
	})

}
