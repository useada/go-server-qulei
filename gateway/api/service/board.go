package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"a.com/go-server/proto/pb"

	"a.com/go-server/gateway/api/base"
	"a.com/go-server/gateway/api/errno"
)

func (s *Service) ListComments(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oid       string `form:"oid" binding:"required"`
		Cid       string `form:"cid"` // cid != "" 拉取二级评论
		PageToken string `form:"page_token"`
	}

	if err := ctx.ShouldBindWith(&args, binding.Query); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	/* for test
	fmt.Println(tracing.GetID(ctx.Request.Context()))
	*/

	res, err := s.Grpc.ListComments(ctx.Request.Context(), &pb.CommListArgs{
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

func (s *Service) GetComment(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		ID  string `form:"id" binding:"required"`
		Oid string `form:"oid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.Query); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := s.Grpc.GetComment(ctx.Request.Context(), &pb.CommGetArgs{
		Id:  args.ID,
		Oid: args.Oid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}

func (s *Service) NewComment(ctx *gin.Context) *base.JSONResponse {
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

	res, err := s.Grpc.NewComment(ctx.Request.Context(), &pb.CommNewArgs{
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

func (s *Service) DelComment(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		ID  string `json:"id" binding:"required"`
		Oid string `json:"oid" binding:"required"`
		Cid string `json:"cid"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := s.Grpc.DelComment(ctx.Request.Context(), &pb.CommDelArgs{
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

func (s *Service) LikeComment(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
		Cid string `json:"cid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := s.Grpc.LikeComment(ctx.Request.Context(), &pb.CommLikeArgs{
		Oid: args.Oid,
		Cid: args.Cid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}

func (s *Service) UnLikeComment(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
		Cid string `json:"cid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := s.Grpc.UnLikeComment(ctx.Request.Context(), &pb.CommUnLikeArgs{
		Oid: args.Oid,
		Cid: args.Cid,
		Uid: "testuid", // TODO
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}

func (s *Service) ListLikes(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oid       string `form:"oid" binding:"required"`
		PageToken string `form:"page_token"`
	}

	if err := ctx.ShouldBindWith(&args, binding.Query); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := s.Grpc.ListLikes(ctx.Request.Context(), &pb.LikeListArgs{
		Oid:       args.Oid,
		PageToken: args.PageToken,
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}

func (s *Service) NewLike(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := s.Grpc.NewLike(ctx.Request.Context(), &pb.LikeNewArgs{
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

func (s *Service) DelLike(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oid string `json:"oid" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := s.Grpc.DelLike(ctx.Request.Context(), &pb.LikeDelArgs{
		Uid: "testuid", // TODO
		Oid: args.Oid,
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}

func (s *Service) GetSummaries(ctx *gin.Context) *base.JSONResponse {
	var args struct {
		Oids []string `json:"oids" binding:"required"`
	}
	if err := ctx.ShouldBindWith(&args, binding.JSON); err != nil {
		return base.ErrorResponse(errno.ARGS_BIND_ERR, err.Error())
	}

	res, err := s.Grpc.GetSummaries(ctx.Request.Context(), &pb.SummaryArgs{
		Uid:  "testuid", // TODO
		Oids: args.Oids,
	})
	if err != nil {
		return base.ErrorResponse(errno.INTERNEL_ERR, err.Error())
	}
	return base.SuccessResponse(res)
}
