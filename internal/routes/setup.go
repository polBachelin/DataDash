package routes

import (
	"dashboard/internal/controllers"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	prefix := "/api/v1"

	r.POST(prefix+"/connect", controllers.SetDatabase)
	r.POST(prefix+"/query", controllers.PostQuery)
}
