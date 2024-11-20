package web

import (
	"io/fs"
	"net/http"
	"strings"
)

type webAssets struct {
	basename string
	assets   fs.FS
}

func NewWebAssets(basename string, assets fs.FS) interface{} {
	return &webAssets{
		basename: basename,
		assets:   assets,
	}
}

func (wa *webAssets) SetupToWeb(app WebApp) error {

	hasIndexFile := false
	entries, err := fs.ReadDir(wa.assets, ".")
	if err != nil {
		return err
	}
	var baseRouter WebRouter
	if wa.basename == "" || wa.basename == "/" {
		baseRouter = app
	} else {
		baseRouter = app.Group(wa.basename)
	}
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() {
			baseRouter.StaticFileFS(name, name, http.FS(wa.assets))
			if !hasIndexFile && name == "index.html" {
				hasIndexFile = true
			}
			continue
		}
		dirFS, err := fs.Sub(wa.assets, name)
		if err != nil {
			return err
		}
		baseRouter.StaticFS(name, http.FS(dirFS))
	}

	if hasIndexFile {
		app.NoRoute(wa.noRoute)
	}

	return nil
}

const (
	MIMEHTML = "text/html"
)

func (wa *webAssets) noRoute(ctx WebContext) {
	accepted := ctx.NegotiateFormat(MIMEHTML) == MIMEHTML
	if accepted && strings.HasPrefix(ctx.Request.URL.Path, wa.basename) {
		ctx.FileFromFS("/", http.FS(wa.assets))
		return
	}
	ctx.Next()
}
