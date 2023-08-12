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

func (x *mongoDatabase) ConnectDatabase(dbData DatabaseInfo) error {
	uri := x.BuildDatabaseUri(dbData)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	x.db = client.Database(dbData.DbName)
	return nil
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

func (x *mongoDatabase) ExecuteQuery(query interface{}) (interface{}, error) {
	//TODO : implement this for mongoDB driver, query needs to be a []bson.M
	return nil, nil
}

func (x *mongoDatabase) QueryResultToJson(result interface{}) ([]map[string]interface{}, error) {
	return nil, nil
}
