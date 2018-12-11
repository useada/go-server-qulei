package main

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	board := r.Group("/v1/")
	board.Use()
	{
		board.GET("/board/comm/list", ResponseWrapper(Board.ListComments))
		board.GET("/board/comm/get", ResponseWrapper(Board.GetComment))
		board.POST("/board/comm/new", ResponseWrapper(Board.NewComment))
		board.POST("/board/comm/del", ResponseWrapper(Board.DelComment))
		board.POST("/board/comm/like", ResponseWrapper(Board.LikeComment))
		board.POST("/board/comm/unlike", ResponseWrapper(Board.UnLikeComment))
		board.GET("/board/like/list", ResponseWrapper(Board.ListLikes))
		board.POST("/board/like/new", ResponseWrapper(Board.NewLike))
		board.POST("/board/like/del", ResponseWrapper(Board.DelLike))
		board.POST("/board/summary/mget", ResponseWrapper(Board.MutiGetSummary)) // 应该在feed接口里被调用
	}

	search := r.Group("/v1/")
	search.Use()
	{
		search.GET("/search/user/name", ResponseWrapper(Search.UsersByName))
		search.GET("/search/user/near", ResponseWrapper(Search.UsersByNear))
	}

	file := r.Group("/v1/")
	file.Use()
	{
		file.POST("/file/upload", ResponseWrapper(File.Upload))
	}

	return r
}
