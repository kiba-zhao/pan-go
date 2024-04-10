package runtime_test

import (
	"pan/runtime"
	"reflect"
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
		simpleModule := runtime.NewSimpleModule(module)
		err := e.Mount(simpleModule)
		assert.Nil(t, err)
	})

	t.Run("test", func(t *testing.T) {
		// module := &TestModule{}
		// simpleModule := runtime.NewSimpleModule(module)
		type_ := reflect.TypeOf(&TestModule{})
		name := type_.Name()
		str := type_.String()
		kind := type_.Kind()
		// fieldnum := type_.NumField()
		// field := type_.Field(0)
		// fieldname := field.Name
		// vfields := reflect.VisibleFields(type_)
		// vexport := vfields[0].IsExported()

		assert.Equal(t, "TestModule", name)
		assert.Equal(t, reflect.Struct, kind)
		assert.Equal(t, "TestModule", str)
		// assert.Equal(t, 1, fieldnum)
		// assert.Equal(t, "module", fieldname)
		// assert.Equal(t, 1, len(vfields))
		// assert.True(t, vexport)
	})
}
