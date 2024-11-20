package web

import "strconv"

const (
	CountHeaderName = "X-Total-Count"
)

func SetCountHeaderForWeb(ctx WebContext, total int64) {
	ctx.Header(CountHeaderName, strconv.FormatInt(total, 10))
	ctx.Header("Access-Control-Expose-Headers'", CountHeaderName)
}
