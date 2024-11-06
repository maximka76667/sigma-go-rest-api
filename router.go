package main

import (
	"github.com/gin-gonic/gin"
	"github.com/maximka76667/sigma-go-rest-api/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/api/pages", controllers.GetPages)
	r.GET("/api/pages/:guid", controllers.GetPageByGUID)
	r.POST("/api/pages", controllers.CreatePage)
	r.PUT("/api/pages/:guid", controllers.UpdatePage)
	r.DELETE("/api/pages/:id", controllers.DeletePage)

	return r
}
