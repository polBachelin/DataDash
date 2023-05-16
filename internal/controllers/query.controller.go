package controllers

import (
	queryService "dashboard/internal/services/query"

	"github.com/gin-gonic/gin"
)

func PostQuery(c *gin.Context) {
	var query queryService.Query
	err := c.BindJSON(&query)
	if err != nil {
		c.JSON(400, "Error in body request")
		return
	}
	queryService.ParseQuery(query)
}
