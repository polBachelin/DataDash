package routes

func Setup(r *gin.Engine) {
	prefix := "/api/v1"

	r.POST(prefix+"/query", controllers.PostQuery)
}
