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

	user := r.Group("/")
	user.Use()
	{
		user.GET("/v1/user/search/name", ResponseWrapper(User.SearchByName))
		user.GET("/v1/user/search/near", ResponseWrapper(User.SearchByNear))
	}

	return r
}
