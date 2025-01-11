package web

import "strconv"

const (
	CountHeaderName = "X-Total-Count"
	ETagHeaderName  = "ETag"
)

func SetCountHeaderForWeb(ctx WebContext, total int64) {
	ctx.Header(CountHeaderName, strconv.FormatInt(total, 10))
	ctx.Header("Access-Control-Expose-Headers'", CountHeaderName)
}

func SetETagHeaderForWeb(ctx WebContext, etag string) {
	ctx.Header(ETagHeaderName, etag)
	ctx.Header("Access-Control-Expose-Headers'", ETagHeaderName)
}
