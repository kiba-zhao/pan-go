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

// NewWebAssets returns a middleware that can be used to serve static files
// from the given fs.FS at the given path.
//
// If basename is empty or "/", the files are served at the root path.
// Otherwise, the given basename is used as the prefix for the path.
//
// For example, if basename is "static", and the assets fs.FS contains a file
// called "foo/bar", the file will be served at "/static/foo/bar".
//
// The middleware will serve an index.html file if one is present in the
// assets fs.FS.
func NewWebAssets(basename string, assets fs.FS) interface{} {
	return &webAssets{
		basename: basename,
		assets:   assets,
	}
}

// Implements WebAppModule interface
//
// SetupToWeb sets up the web application to serve the static files from the given
// fs.FS at the given path.
//
// If basename is empty or "/", the files are served at the root path.
// Otherwise, the given basename is used as the prefix for the path.
//
// The middleware will serve an index.html file if one is present in the
// assets fs.FS. If an index.html file is present, it will also be served
// when the user visits the root path.
//
// For example, if basename is "static", and the assets fs.FS contains a file
// called "foo/bar", the file will be served at "/static/foo/bar".
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

// noRoute handles requests that do not match any defined routes.
// If the request accepts HTML and the URL path starts with the specified basename,
// it serves the index.html file from the root of the web assets file system.
// Otherwise, it proceeds to the next middleware or route handler.

func (wa *webAssets) noRoute(ctx WebContext) {
	accepted := ctx.NegotiateFormat(MIMEHTML) == MIMEHTML
	if accepted && strings.HasPrefix(ctx.Request.URL.Path, wa.basename) {
		ctx.FileFromFS("/", http.FS(wa.assets))
		return
	}
	ctx.Next()
}
