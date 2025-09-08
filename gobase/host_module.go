//go:build !(android || ios) || host

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
