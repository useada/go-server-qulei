package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"a.com/go-server/gclient"
	"a.com/go-server/proto/pb"
)

type BoardHandler struct {
}

var Board *BoardHandler

func (b *BoardHandler) ListComments(ctx *gin.Context) *JSONResponse {
	var args struct {
		Oid       string `form:"oid" binding:"required"`
		Cid       string `form:"cid"` // cid != "" 拉取二级评论
		PageToken string `form:"page_token"`
	}

	if err := ctx.ShouldBindWith(&args, binding.Query); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	/* for test
	fmt.Println(tracing.GetID(ctx.Request.Context()))
	*/

	res, err := gclient.Board.ListComments(ctx.Request.Context(), &pb.CommListArgs{
		Oid:       args.Oid,
		Cid:       args.Cid,
		Uid:       "testuid", // TODO
		PageToken: args.PageToken,
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) GetComment(ctx *gin.Context) *JSONResponse {
	var args struct {
		ID  string `form:"id" binding:"required"`
		Oid string `form:"oid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.Query); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.GetComment(ctx.Request.Context(), &pb.CommGetArgs{
		Id:  args.ID,
		Oid: args.Oid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) NewComment(ctx *gin.Context) *JSONResponse {
	var args struct {
		Oid      string `json:"oid" binding:"required"`
		Cid      string `json:"cid"`
		IsRepost bool   `json:"is_repost"`
		Content  string `json:"content" binding:"lte=2000"`
		ImgID    string `json:"img_id"`
		ImgEx    string `json:"img_ex"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.NewComment(ctx.Request.Context(), &pb.CommNewArgs{
		Oid: args.Oid,
		Cid: args.Cid,
		Author: &pb.UserBaseInfo{
			Uid: "testuid", // TODO
		},
		IsRepost: args.IsRepost,
		ImgId:    args.ImgID,
		ImgEx:    args.ImgEx,
		Content:  args.Content,
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) DelComment(ctx *gin.Context) *JSONResponse {
	var args struct {
		ID  string `json:"id" binding:"required"`
		Oid string `json:"oid" binding:"required"`
		Cid string `json:"cid"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.DelComment(ctx.Request.Context(), &pb.CommDelArgs{
		Id:  args.ID,
		Oid: args.Oid,
		Cid: args.Cid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) LikeComment(ctx *gin.Context) *JSONResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
		Cid string `json:"cid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.LikeComment(ctx.Request.Context(), &pb.CommLikeArgs{
		Oid: args.Oid,
		Cid: args.Cid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) UnLikeComment(ctx *gin.Context) *JSONResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
		Cid string `json:"cid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.UnLikeComment(ctx.Request.Context(), &pb.CommUnLikeArgs{
		Oid: args.Oid,
		Cid: args.Cid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) ListLikes(ctx *gin.Context) *JSONResponse {
	var args struct {
		Oid       string `form:"oid" binding:"required"`
		PageToken string `form:"page_token"`
	}

	if err := ctx.ShouldBindWith(&args, binding.Query); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.ListLikes(ctx.Request.Context(), &pb.LikeListArgs{
		Oid:       args.Oid,
		PageToken: args.PageToken,
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) NewLike(ctx *gin.Context) *JSONResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.NewLike(ctx.Request.Context(), &pb.LikeNewArgs{
		Author: &pb.UserBaseInfo{
			Uid: "testuid", // TODO
		},
		Oid: args.Oid,
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) DelLike(ctx *gin.Context) *JSONResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.DelLike(ctx.Request.Context(), &pb.LikeDelArgs{
		Uid: "testuid", // TODO
		Oid: args.Oid,
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) MutiGetSummary(ctx *gin.Context) *JSONResponse {
	var args struct {
		Oids []string `json:"oids" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.MutiGetSummary(ctx.Request.Context(), &pb.BoardSummaryArgs{
		Uid:  "testuid", // TODO
		Oids: args.Oids,
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}