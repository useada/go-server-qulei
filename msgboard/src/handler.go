package main

import (
	"google.golang.org/grpc"

	proto "a.com/go-server/proto/pb/msgboard"
)

func RegisterHandler(svr *grpc.Server) {
	proto.RegisterMsgBoardServer(svr, &SvrHandler{})
}

type SvrHandler struct{}
