package runtime_test

import (
	"pan/internal/runtime"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleModule(t *testing.T) {

	setup := func() (engine *runtime.Engine) {
		engine, _ = runtime.New()
		return
	}

	t.Run("simple module", func(t *testing.T) {
		e := setup()

		module := &TestModule{}
		simpleModule := runtime.NewModule(module)
		err := e.Mount(simpleModule)
		assert.Nil(t, err)
	})

	t.Run("smaple test", func(t *testing.T) {
		var iface TestInterface

		iface = &sampleTestModule{}
		assert.Equal(t, iface.Hello(), "hello")

		realType := reflect.TypeFor[*sampleTestModule]()
		ifaceType := reflect.TypeOf(iface)
		print(ifaceType == realType)
	})

}

type TestInterface interface {
	Hello() string
}

type sampleTestModule struct {
}

func (t *sampleTestModule) Hello() string {
	return "hello"
}

type sampleTestModuleNext struct {
	sampleTestModule
}
