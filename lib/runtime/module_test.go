package runtime_test

import (
	"net/url"
	"pan/lib/runtime"
	"path/filepath"
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

	t.Run("test", func(t *testing.T) {
		// urlString := "https://raw.githubusercontent.com/kiba-zhao/pan-go/refs/heads/dev/README.zh-CN.md"
		// urlString := "c:\\Users\\kiba-zhao"
		// urlString := "/tmp/test-txt"
		urlString := "file:///C:/Users/example/Documents/file.txt"
		vol, err := url.ParseRequestURI(urlString)

		assert.Nil(t, err)
		assert.NotNil(t, vol)

		fp := filepath.FromSlash(vol.Path)
		assert.NotNil(t, fp)
	})
}
