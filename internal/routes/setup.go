package routes

import (
	"dashboard/internal/controllers"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.POST("/connect", controllers.SetDatabase)
	r.POST("/query", controllers.PostQuery)
}
