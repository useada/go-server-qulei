package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"a.com/go-server/proto/pb"

	"a.com/go-server/gateway/api/internal/base"
	"a.com/go-server/gateway/api/internal/errno"
)

func (h *Handler) ListComments(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oid       string `form:"oid" binding:"required"`
		Cid       string `form:"cid"` // cid != "" 拉取二级评论
		PageToken string `form:"page_token"`
	}

	if err := ctx.ShouldBindWith(&args, binding.Query); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	fmt.Println("--------", h.TraceId(ctx))

	res, err := h.Grpc.ListComments(ctx.Request.Context(), &pb.CommListArgs{
		Oid:       args.Oid,
		Cid:       args.Cid,
		Uid:       "testuid", // TODO
		PageToken: args.PageToken,
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}

func (h *Handler) GetComment(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		ID  string `form:"id" binding:"required"`
		Oid string `form:"oid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.Query); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := h.Grpc.GetComment(ctx.Request.Context(), &pb.CommGetArgs{
		Id:  args.ID,
		Oid: args.Oid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}

func (h *Handler) NewComment(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oid      string `json:"oid" binding:"required"`
		Cid      string `json:"cid"`
		IsRepost bool   `json:"is_repost"`
		Content  string `json:"content" binding:"lte=2000"`
		ImgID    string `json:"img_id"`
		ImgEx    string `json:"img_ex"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := h.Grpc.NewComment(ctx.Request.Context(), &pb.CommNewArgs{
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
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}

func (h *Handler) DelComment(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		ID  string `json:"id" binding:"required"`
		Oid string `json:"oid" binding:"required"`
		Cid string `json:"cid"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := h.Grpc.DelComment(ctx.Request.Context(), &pb.CommDelArgs{
		Id:  args.ID,
		Oid: args.Oid,
		Cid: args.Cid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}

func (h *Handler) LikeComment(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
		Cid string `json:"cid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := h.Grpc.LikeComment(ctx.Request.Context(), &pb.CommLikeArgs{
		Oid: args.Oid,
		Cid: args.Cid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}

func (h *Handler) UnLikeComment(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
		Cid string `json:"cid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := h.Grpc.UnLikeComment(ctx.Request.Context(), &pb.CommUnLikeArgs{
		Oid: args.Oid,
		Cid: args.Cid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}

func (h *Handler) ListLikes(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oid       string `form:"oid" binding:"required"`
		PageToken string `form:"page_token"`
	}

	if err := ctx.ShouldBindWith(&args, binding.Query); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := h.Grpc.ListLikes(ctx.Request.Context(), &pb.LikeListArgs{
		Oid:       args.Oid,
		PageToken: args.PageToken,
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}

func (h *Handler) NewLike(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := h.Grpc.NewLike(ctx.Request.Context(), &pb.LikeNewArgs{
		Author: &pb.UserBaseInfo{
			Uid: "testuid", // TODO
		},
		Oid: args.Oid,
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}

func (h *Handler) DelLike(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := h.Grpc.DelLike(ctx.Request.Context(), &pb.LikeDelArgs{
		Uid: "testuid", // TODO
		Oid: args.Oid,
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}

func (h *Handler) GetSummaries(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oids []string `json:"oids" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := h.Grpc.GetSummaries(ctx.Request.Context(), &pb.SummaryArgs{
		Uid:  "testuid", // TODO
		Oids: args.Oids,
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}
