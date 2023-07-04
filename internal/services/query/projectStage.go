package query

import "go.mongodb.org/mongo-driver/bson"

func generateProjectStage(dimensions, measures []string) bson.M {
	projectStage := bson.M{
		"$project": bson.M{"_id": 0},
	}

	for _, dimension := range dimensions {
		memberName := getMemberName(dimension)
		projectStage["$project"].(bson.M)[memberName] = "$_id." + memberName
	}

	for _, measure := range measures {
		projectStage["$project"].(bson.M)[getMemberName(measure)] = 1
	}
	return projectStage
}
