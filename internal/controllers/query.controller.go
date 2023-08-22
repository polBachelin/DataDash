package controllers

import (
	"dashboard/internal/database"
	"dashboard/internal/services/block"
	queryService "dashboard/internal/services/query"
	"fmt"

	"github.com/gin-gonic/gin"
)

type QueryResponse struct {
	Query       queryService.Query       `json:"query"`
	Data        []map[string]interface{} `json:"data"`
	Annotations queryService.Annotations `json:"annotation"`
}

func PostQuery(c *gin.Context) {
	var query queryService.Query

	if database.GetPostgresDatabase() == nil {
		c.JSON(400, "Error please connect to database first")
		return
	}
	err := c.BindJSON(&query)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Error in body request")
		return
	}
	b := block.GetInstance().Blocks
	service := queryService.NewQueryService(query, database.GetCurrentDatabase(), block.NewGraph(b))
	data, err := service.ParseQuery()
	if err != nil {
		c.JSON(500, "Internal error")
		return
	}
	queryResponse := &QueryResponse{Query: query, Data: data, Annotations: service.CreateAnnotations()}
	c.JSON(200, queryResponse)
}
