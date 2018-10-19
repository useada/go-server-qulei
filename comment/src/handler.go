package main

import (
	"google.golang.org/grpc"

	proto "a.com/go-server/proto/pb/comment"
)

func RegisterHandler(svr *grpc.Server) {
	proto.RegisterCommentServer(svr, &SvrHandler{})
}

type SvrHandler struct{}
