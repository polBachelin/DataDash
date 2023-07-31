package controllers

import (
	"dashboard/internal/database"
	queryService "dashboard/internal/services/query"
	"fmt"

	"github.com/gin-gonic/gin"
)

func PostQuery(c *gin.Context) {
	var query queryService.Query

	if database.GetMongoDatabase() == nil {
		c.JSON(400, "Error please connect to database first")
		return
	}
	err := c.BindJSON(&query)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Error in body request")
		return
	}
	service := queryService.NewQueryService(query, database.GetCurrentDatabase())
	res, err := service.ParseQuery()
	if err != nil {
		c.JSON(500, "Internal error")
	}
	c.JSON(200, res)
}
