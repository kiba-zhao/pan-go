//go:build !(android || ios)

package gobase

import (
	"pan/lib/runtime"
	libWeb "pan/lib/web"
	"pan/web"
)

func NewHostModule() interface{} {
	return runtime.NewModule(
		libWeb.New(),
		libWeb.NewWebAssets("/", web.AssetFS()),
	)
}
