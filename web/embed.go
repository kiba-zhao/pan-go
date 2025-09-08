package web

import (
	"embed"
	"io/fs"
	"pan/lib/log"
)

//go:embed dist
var assetFS embed.FS

func AssetFS() fs.FS {
	logger := log.Default()
	assetFS_, err := fs.Sub(assetFS, "dist")
	if err != nil {
		logger.Error("web", "AssetFS Error:"+err.Error())
		return nil
	}
	return assetFS_
}
