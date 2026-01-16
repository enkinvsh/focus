package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	api.Use(AuthMiddleware())
	{
		api.GET("/tasks", GetTasks)
		api.POST("/tasks", CreateTask)
		api.PATCH("/tasks/:id", UpdateTask)
		api.DELETE("/tasks/:id", DeleteTask)

		api.GET("/user/preferences", GetPreferences)
		api.PATCH("/user/preferences", UpdatePreferences)
	}
}
