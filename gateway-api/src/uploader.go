package main

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"a.com/go-server/common/crypto"
	"a.com/go-server/gclient"
	"a.com/go-server/proto/pb"
)

type FileHandler struct {
}

var File *FileHandler

func (f *FileHandler) Upload(ctx *gin.Context) *JsonResponse {
	var args struct {
		Sha  string `form:"sha" binding:"required"`
		Ex   string `form:"ex" binding:"required"`
		Type int    `form:"type" binding:"min=0,max=2"` // 0 image 1 avatar 2 file
	}
	if err := ctx.ShouldBindWith(&args, binding.Query); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	if res, err := gclient.Uploader.Query(&pb.FileQueryArgs{
		Id: args.Sha,
	}); err == nil {
		return SuccessResponse(res)
	}

	bin, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return ErrorResponse(ARGS_READ_ERR, err.Error())
	}

	id, err := crypto.Sha1(string(bin))
	if err != nil {
		Log.Error("sha1 file data err:%v", err)
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}

	res, err := gclient.Uploader.Upload(&pb.FileUploadArgs{
		Id:    id,
		Ex:    args.Ex,
		Typef: pb.FileType(args.Type),
		Data:  bin,
	})
	if err != nil {
		Log.Error("upload file err:%v", err)
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}
