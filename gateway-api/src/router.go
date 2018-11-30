package main

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	file := r.Group("/")
	file.Use()
	{
		file.POST("/v1/file/upload", ResponseWrapper(File.Upload))
	}

	search := r.Group("/")
	search.Use()
	{
		search.GET("/v1/search/user/name", ResponseWrapper(Search.UsersByName))
		search.GET("/v1/search/user/near", ResponseWrapper(Search.UsersByNear))
	}

	return r
}
