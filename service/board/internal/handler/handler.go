package handler

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"a.com/go-server/proto/pb"

	"a.com/go-server/service/board/internal/cache"
	"a.com/go-server/service/board/internal/store"
)

const (
	PAGE_COUNT = 20

	COUNT_COMM_CACHE = 200
)

func RegisterHandler(svr *grpc.Server, kv cache.Cache, db store.Store, log *zap.SugaredLogger) {
	pb.RegisterBoardServer(svr, &SvrHandler{
		Cache: kv,
		Store: db,
		Log:   log,
	})
}

type SvrHandler struct {
	Cache cache.Cache
	Store store.Store
	Log   *zap.SugaredLogger
}
