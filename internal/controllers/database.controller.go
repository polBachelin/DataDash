package controllers

import (
	"dashboard/internal/database"
	"dashboard/pkg/utils"
	"fmt"

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
	postgres := database.GetPostgresDatabase()
	res := postgres.ConnectDatabase(dbData)
	if res != nil {
		c.JSON(500, "Error connecting to database"+res.Error())
		return
	}
	database.SetPostgresDatabase(postgres)
	c.JSON(200, "Successfully connected to database")
}

func SetDatabaseFromEnv() error {
	var dbData database.DatabaseInfo

	dbData.DbHost = utils.GetEnvVar("DB_HOST", "0.0.0.0")
	dbData.DbPort = utils.GetEnvVar("DB_POST", "5432")
	dbData.DbName = utils.GetEnvVar("DB_NAME", "postgres")
	dbData.DbUsername = utils.GetEnvVar("DB_USER", "postgres")
	dbData.DbPass = utils.GetEnvVar("DB_PASS", "postgres")
	postgres := database.GetPostgresDatabase()
	res := postgres.ConnectDatabase(dbData)
	if res != nil {
		return fmt.Errorf("Could not connect to database %s: %v", dbData.DbName, res)
	}
	database.SetPostgresDatabase(postgres)
	return nil
}
