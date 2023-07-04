package query

import "go.mongodb.org/mongo-driver/bson"

func generateLimitStage(limit int) bson.M {
	return bson.M{"$limit": limit}
}
