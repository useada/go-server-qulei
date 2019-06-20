package main

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"a.com/go-server/common/profiler"
	"a.com/go-server/common/tracing"
)

func Router() *gin.Engine {
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
		board.GET("/comm/list", ResponseWrapper(Board.ListComments))
		board.GET("/comm/get", ResponseWrapper(Board.GetComment))
		board.POST("/comm/new", ResponseWrapper(Board.NewComment))
		board.POST("/comm/del", ResponseWrapper(Board.DelComment))
		board.POST("/comm/like", ResponseWrapper(Board.LikeComment))
		board.POST("/comm/unlike", ResponseWrapper(Board.UnLikeComment))
		board.GET("/like/list", ResponseWrapper(Board.ListLikes))
		board.POST("/like/new", ResponseWrapper(Board.NewLike))
		board.POST("/like/del", ResponseWrapper(Board.DelLike))
		board.POST("/summary/mget", ResponseWrapper(Board.MutiGetSummary)) // 应该在feed接口里被调用
	}

	// 搜索
	search := router.Group("/v1/search/")
	{
		search.GET("/user/name", ResponseWrapper(Search.UsersByName))
		search.GET("/user/near", ResponseWrapper(Search.UsersByNear))
	}

	// 文件上传
	file := router.Group("/v1/file")
	{
		file.POST("/upload", ResponseWrapper(File.Upload))
	}

	return router
}
