package service

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"a.com/go-server/common/crypto"
	"a.com/go-server/gateway/api/base"
	"a.com/go-server/gateway/api/errno"
	"a.com/go-server/proto/pb"
)

func (s *Service) Upload(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Sha  string `form:"sha" binding:"required"`
		Ex   string `form:"ex" binding:"required"`
		Type int    `form:"type" binding:"min=0,max=2"` // 0 image 1 avatar 2 file
	}
	if err := ctx.ShouldBindWith(&args, binding.Query); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	if res, err := s.Grpc.Query(ctx.Request.Context(), &pb.FileQueryArgs{
		Id: args.Sha,
	}); err == nil {
		return base.SuccessResponse(res)
	}

	bin, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return base.ErrorResponse(errno.ARGS_READ_ERR, err.Error())
	}

	id, err := crypto.Sha1(string(bin))
	if err != nil {
		s.Log.Error("sha1 file data err:%v", err)
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}

	res, err := s.Grpc.Upload(ctx.Request.Context(), &pb.FileUploadArgs{
		Id:   id,
		Ex:   args.Ex,
		Type: pb.TYPE(args.Type),
		Data: bin,
	})
	if err != nil {
		s.Log.Error("upload file err:%v", err)
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}
