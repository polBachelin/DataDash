package main

import (
	"dashboard/internal/controllers"
	"dashboard/internal/routes"
	"dashboard/pkg/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, access-control-allow-origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(CORS())
	err := controllers.SetDatabaseFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	routes.Setup(r)
	r.Run(":" + utils.GetEnvVar("API_PORT", "8080"))
}
