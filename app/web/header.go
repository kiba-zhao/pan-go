package web

import "strconv"

const (
	CountHeaderName = "X-Total-Count"
	ETagHeaderName  = "ETag"
)

// SetCountHeaderForWeb sets the HTTP header "X-Total-Count" to the given value on the given context.
// This is used to expose the total count of items in a paginated response to the client.
// The "Access-Control-Expose-Headers" header is also set to allow the client to access the "X-Total-Count" header.
func SetCountHeaderForWeb(ctx WebContext, total int64) {
	ctx.Header(CountHeaderName, strconv.FormatInt(total, 10))
	ctx.Header("Access-Control-Expose-Headers'", CountHeaderName)
}

// SetETagHeaderForWeb sets the HTTP header "ETag" to the given value on the given context.
// This is used to expose the ETag of a resource to the client.
// The "Access-Control-Expose-Headers" header is also set to allow the client to access the "ETag" header.
func SetETagHeaderForWeb(ctx WebContext, etag string) {
	ctx.Header(ETagHeaderName, etag)
	ctx.Header("Access-Control-Expose-Headers'", ETagHeaderName)
}
