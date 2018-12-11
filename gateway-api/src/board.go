package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"a.com/go-server/common/page"
	"a.com/go-server/gclient"
	"a.com/go-server/proto/ct"
	"a.com/go-server/proto/pb"
)

type BoardHandler struct {
}

var Board *BoardHandler

const (
	BOARD_PAGE_COUNT = 20
)

func (b *BoardHandler) ListComments(ctx *gin.Context) *JsonResponse {
	var args struct {
		Oid       string `form:"oid" binding:"required"`
		Cid       string `form:"cid"` // cid != "" 拉取二级评论
		PageToken string `form:"page_token"`
	}
	if err := ctx.ShouldBindWith(&args, binding.Query); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}
	if args.PageToken == "" {
		args.PageToken, _ = page.Default(ct.TIME_INF_MAX, BOARD_PAGE_COUNT)
	}

	res, err := gclient.Board.ListComments(&pb.CommListArgs{
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

func (b *BoardHandler) GetComment(ctx *gin.Context) *JsonResponse {
	var args struct {
		Id  string `form:"id" binding:"required"`
		Oid string `form:"oid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.Query); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.GetComment(&pb.CommGetArgs{
		Id:  args.Id,
		Oid: args.Oid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) NewComment(ctx *gin.Context) *JsonResponse {
	var args struct {
		Oid      string `json:"oid" binding:"required"`
		Cid      string `json:"cid"`
		IsRepost bool   `json:"is_repost"`
		Content  string `json:"content" binding:"lte=2000"`
		ImgId    string `json:"img_id"`
		ImgEx    string `json:"img_ex"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.NewComment(&pb.CommNewArgs{
		Oid: args.Oid,
		Cid: args.Cid,
		Author: &pb.UserBaseInfo{
			Uid: "testuid", // TODO
		},
		IsRepost: args.IsRepost,
		ImgId:    args.ImgId,
		ImgEx:    args.ImgEx,
		Content:  args.Content,
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) DelComment(ctx *gin.Context) *JsonResponse {
	var args struct {
		Id  string `json:"id" binding:"required"`
		Oid string `json:"oid" binding:"required"`
		Cid string `json:"cid"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.DelComment(&pb.CommDelArgs{
		Id:  args.Id,
		Oid: args.Oid,
		Cid: args.Cid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) LikeComment(ctx *gin.Context) *JsonResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
		Cid string `json:"cid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.LikeComment(&pb.CommLikeArgs{
		Oid: args.Oid,
		Cid: args.Cid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) UnLikeComment(ctx *gin.Context) *JsonResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
		Cid string `json:"cid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.UnLikeComment(&pb.CommUnLikeArgs{
		Oid: args.Oid,
		Cid: args.Cid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) ListLikes(ctx *gin.Context) *JsonResponse {
	var args struct {
		Oid       string `form:"oid" binding:"required"`
		PageToken string `form:"page_token"`
	}
	if err := ctx.ShouldBindWith(&args, binding.Query); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}
	if args.PageToken == "" {
		args.PageToken, _ = page.Default(ct.TIME_INF_MAX, BOARD_PAGE_COUNT)
	}

	res, err := gclient.Board.ListLikes(&pb.LikeListArgs{
		Oid:       args.Oid,
		PageToken: args.PageToken,
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) NewLike(ctx *gin.Context) *JsonResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.NewLike(&pb.LikeNewArgs{
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

func (b *BoardHandler) DelLike(ctx *gin.Context) *JsonResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.DelLike(&pb.LikeDelArgs{
		Uid: "testuid", // TODO
		Oid: args.Oid,
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}

func (b *BoardHandler) MutiGetSummary(ctx *gin.Context) *JsonResponse {
	var args struct {
		Oids []string `json:"oids" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return ErrorResponse(ARGS_BIND_ERR, err.Error())
	}

	res, err := gclient.Board.MutiGetSummary(&pb.BoardSummaryArgs{
		Uid:  "testuid", // TODO
		Oids: args.Oids,
	})
	if err != nil {
		return ErrorResponse(INTERNEL_ERR, err.Error())
	}
	return SuccessResponse(res)
}
