//go:build android || ios

package gomobile

import "io"

type GoMobileReadStream = io.Reader

type GoMobileStream = io.ReadCloser
