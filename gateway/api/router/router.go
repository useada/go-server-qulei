package router

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"a.com/go-server/common/profiler"
	"a.com/go-server/common/tracing"

	"a.com/go-server/gateway/api/base"
	"a.com/go-server/gateway/api/service"
)

func BindRouter(s *service.Service) *gin.Engine {
	router := gin.Default()

	// middleware
	{
		router.Use(gin.Recovery())
		router.Use(tracing.Trace(opentracing.GlobalTracer()))
	}

	// health check
	{
		router.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })
	}

	// pprof
	{
		profiler.Prof(router)
	}

	// 评论点赞
	board := router.Group("/v1/board/")
	{
		board.GET("/comm/list", base.ResponseWrapper(s.ListComments))
		board.GET("/comm/get", base.ResponseWrapper(s.GetComment))
		board.POST("/comm/new", base.ResponseWrapper(s.NewComment))
		board.POST("/comm/del", base.ResponseWrapper(s.DelComment))
		board.POST("/comm/like", base.ResponseWrapper(s.LikeComment))
		board.POST("/comm/unlike", base.ResponseWrapper(s.UnLikeComment))
		board.GET("/like/list", base.ResponseWrapper(s.ListLikes))
		board.POST("/like/new", base.ResponseWrapper(s.NewLike))
		board.POST("/like/del", base.ResponseWrapper(s.DelLike))
		board.POST("/summary/gets", base.ResponseWrapper(s.GetSummaries)) // 应该在feed接口里被调用
	}

	// 搜索
	search := router.Group("/v1/search/")
	{
		search.GET("/user/name", base.ResponseWrapper(s.UsersByName))
		search.GET("/user/near", base.ResponseWrapper(s.UsersByNear))
	}

	// 文件上传
	file := router.Group("/v1/file")
	{
		file.POST("/upload", base.ResponseWrapper(s.Upload))
	}

	return router
}
