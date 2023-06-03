package controllers

import (
	"dashboard/internal/database"

	"github.com/gin-gonic/gin"
)

func SetDatabase(c *gin.Context) {
	var dbData database.DatabaseInfo

	err := c.BindJSON(&dbData)
	if err != nil || dbData.DbHost == "" ||
		dbData.DbPass == "" ||
		dbData.DbPort == "" ||
		dbData.DbUsername == "" ||
		dbData.DbName == "" {
		c.JSON(400, "Error in body request")
		return
	}
	res := database.ConnectDatabase(dbData)
	if res == nil {
		c.JSON(500, "Error connecting to database")
	}
	c.JSON(200, "Successfully connected to database")
}
