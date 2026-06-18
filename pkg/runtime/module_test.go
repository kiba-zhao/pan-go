package runtime_test

import (
	"math"
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

	t.Run("smaple test", func(t *testing.T) {
		prev := uint(0)
		next := uint(math.MaxUint32)

		results := int(prev - next)
		assert.Equal(t, results, -1)
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
