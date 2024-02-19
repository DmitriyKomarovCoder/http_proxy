package http

import "github.com/gin-gonic/gin"

func InitRouter(handler *Handler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")

	api.GET("/requests", handler.AllRequest)
	api.GET("/request/:id", handler.GetRequest)
	api.GET("/repeat/:id", handler.Repeat)
	api.GET("/scan/:id", handler.Scan)

	return router
}
