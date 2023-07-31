package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDatabase struct {
	db *mongo.Database
}

func (x *mongoDatabase) BuildDatabaseUri(dbData DatabaseInfo) string {
	return "mongodb://" + dbData.DbUsername + ":" + dbData.DbPass + "@" + dbData.DbHost + ":" + dbData.DbPort
}

func (x *mongoDatabase) ConnectDatabase(dbData DatabaseInfo) *mongo.Database {
	uri := x.BuildDatabaseUri(dbData)
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
	x.db = client.Database(dbData.DbName)
	return x.db
}

func (x *mongoDatabase) GetDatabaseConnection() *mongo.Database {
	if x.db != nil {
		return x.db
	}
	fmt.Println("No database connection")
	return nil
}

func (x *mongoDatabase) GetCollection(name string) *mongo.Collection {
	return x.GetDatabaseConnection().Collection(name)
}
