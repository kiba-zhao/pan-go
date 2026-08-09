package runtime_test

import (
	"pan/pkg/runtime"
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
