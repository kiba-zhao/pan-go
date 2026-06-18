package net

import (
	"context"
	"errors"
	"io"
	"pan/pkg/servlet"
)

type PeerServer interface {
	Serve(ctx context.Context, conn PeerConn) error
}

var ErrPeerServerInvalidApplet = errors.New("peer.PeerServer Error: Invalid Applet")
var ErrPeerServerNotFound = errors.New("peer.PeerServer Error: Not Found")

type PeerServletContext = *servlet.Context
type PeerServlet = *servlet.Servlet[PeerServletContext]
type PeerServletNext = servlet.Next
type PeerServletRouter = servlet.HandleGroup[PeerServletContext]
type PeerServletScope = servlet.RequestName

const (
	CodeOK            = servlet.CodeOK
	CodeInternalError = 500
	CodeNotFound      = 404
	CodeBadRequest    = 400
	CodeForbidden     = 403
)

var (
	PeerIDSessionKey   = []byte("PeerID")
	PeerConnSessionKey = []byte("PeerConn")
)

func servePeerStream(peerServlet PeerServlet, conn PeerConn, stream PeerStream, target PeerID) error {
	defer stream.Close()

	var err error
	ctx := servlet.NewContext()
	if peerServlet == nil {
		err = ErrPeerServerInvalidApplet
	} else {
		err = servlet.UnmarshalRequest(stream, ctx.Request())
	}

	if err == nil {
		ctx.Set(PeerIDSessionKey, target)
		ctx.Set(PeerConnSessionKey, conn)
		err = peerServlet.Run(ctx, nil)
		defer ctx.Close()
	}

	if err != nil {
		errCode := CodeInternalError
		if codeErr, ok := err.(CodeError); ok {
			errCode = codeErr.Code()
		}
		ctx.ThrowError(errCode, err)
	}

	if ctx.Code() < 0 {
		ctx.ThrowError(CodeNotFound, ErrPeerServerNotFound)
	}

	reader := servlet.MarshalResponse(&ctx.Response)
	_, resErr := io.Copy(stream, reader)

	if err == nil && resErr != nil {
		err = resErr
	}

	return err
}
