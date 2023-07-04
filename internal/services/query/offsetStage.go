package query

import "go.mongodb.org/mongo-driver/bson"

func generateOffsetStage(offset int) bson.M {
	return bson.M{"$skip": offset}
}
