package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseInfo struct {
	DbHost     string `json:"db_host"`
	DbPort     string `json:"db_port"`
	DbUsername string `json:"db_username"`
	DbPass     string `json:"db_pass"`
	DbName     string `json:"db_name"`
}

type database struct {
	db *mongo.Database
}

var databaseConnection = &database{
	db: nil,
}

func BuildDatabaseUri(dbData DatabaseInfo) string {
	return "mongodb://" + dbData.DbUsername + ":" + dbData.DbPass + "@" + dbData.DbHost + ":" + dbData.DbPort
}

func ConnectDatabase(dbData DatabaseInfo) *mongo.Database {
	uri := BuildDatabaseUri(dbData)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	databaseConnection.db = client.Database(dbData.DbName)
	return databaseConnection.db
}

func GetDatabaseConnection() *mongo.Database {
	if databaseConnection.db != nil {
		return databaseConnection.db
	}
	fmt.Println("No database connection")
	return nil
}

func GetCollection(name string) *mongo.Collection {
	return GetDatabaseConnection().Collection(name)
}
