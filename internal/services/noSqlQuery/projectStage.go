package noSqlQuery

import (
	"dashboard/internal/services/query"

	"go.mongodb.org/mongo-driver/bson"
)

func generateProjectStage(dimensions, measures []string) bson.M {
	projectStage := bson.M{
		"$project": bson.M{"_id": 0},
	}

	for _, dimension := range dimensions {
		memberName := query.GetMemberName(dimension)
		projectStage["$project"].(bson.M)[memberName] = "$_id." + memberName
	}

	for _, measure := range measures {
		projectStage["$project"].(bson.M)[query.GetMemberName(measure)] = 1
	}
	return projectStage
}
