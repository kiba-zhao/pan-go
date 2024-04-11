package runtime_test

import (
	"pan/runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleModule(t *testing.T) {

	setup := func() (engine *runtime.Engine) {
		engine = runtime.New()
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
