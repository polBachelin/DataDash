package query

import (
	"context"
	"dashboard/internal/database"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func executeStages(stages []bson.M, collectionName string) []bson.M {
	collection := database.GetCollection(collectionName)
	res, err := collection.Aggregate(context.TODO(), stages)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	document := []bson.M{}
	err = res.All(context.TODO(), &document)
	if err != nil {
		return nil
	}
	return document
}

// calls the given function with the array of measurements to return the MongoDB pipeline
func BuildAStage[E any](s []E, f func([]E) ([]bson.M, error)) ([]bson.M, error) {
	stage, err := f(s)
	if err != nil {
		return []bson.M{}, err
	}
	return stage, nil
}
