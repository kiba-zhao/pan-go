//go:build !(android || ios)

package gobase

import (
	"pan/pkg/runtime"
	libWeb "pan/pkg/web"
	"pan/web"
)

func NewHostModule() interface{} {
	return runtime.NewModule(
		libWeb.New(),
		libWeb.NewWebAssets("/", web.AssetFS()),
	)
}
