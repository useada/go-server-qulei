package router

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"a.com/go-server/common/profiler"
	"a.com/go-server/common/tracing"

	"a.com/go-server/gateway/api/internal/base"
	"a.com/go-server/gateway/api/internal/handler"
)

func BindRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	// 中间件 健康检查
	prepare(router)

	// 评论
	board("/v1/board/", router, h)

	// 搜索
	search("/v1/search/", router, h)

	// 文件上传
	upload("/v1/upload/", router, h)

	return router
}

func board(uri string, router *gin.Engine, h *handler.Handler) {
	board := router.Group(uri)

	board.GET("/comm/list", base.ResponseWrapper(h.ListComments))
	board.GET("/comm/get", base.ResponseWrapper(h.GetComment))
	board.POST("/comm/new", base.ResponseWrapper(h.NewComment))
	board.POST("/comm/del", base.ResponseWrapper(h.DelComment))

	board.POST("/comm/like", base.ResponseWrapper(h.LikeComment))
	board.POST("/comm/unlike", base.ResponseWrapper(h.UnLikeComment))

	board.GET("/like/list", base.ResponseWrapper(h.ListLikes))
	board.POST("/like/new", base.ResponseWrapper(h.NewLike))
	board.POST("/like/del", base.ResponseWrapper(h.DelLike))

	board.POST("/summary/gets", base.ResponseWrapper(h.GetSummaries)) // 应该在feed接口里被调用
}

func search(uri string, router *gin.Engine, h *handler.Handler) {
	search := router.Group(uri)

	search.GET("/user/name", base.ResponseWrapper(h.UsersByName))
	search.GET("/user/near", base.ResponseWrapper(h.UsersByNear))
}

func upload(uri string, router *gin.Engine, h *handler.Handler) {
	upload := router.Group(uri)

	upload.POST("/upload", base.ResponseWrapper(h.Upload))
}

func prepare(router *gin.Engine) {
	// middleware
	router.Use(gin.Recovery())
	router.Use(tracing.Trace(opentracing.GlobalTracer()))

	// health check
	router.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	// pprof
	profiler.Prof(router)
}
