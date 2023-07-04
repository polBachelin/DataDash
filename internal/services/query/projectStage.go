package query

import "go.mongodb.org/mongo-driver/bson"

func generateProjectStage(dimensions, measures []string) bson.M {
	projectStage := bson.M{
		"$project": bson.M{"_id": 0},
	}

	for _, dimension := range dimensions {
		projectStage["$project"].(bson.M)[getMemberName(dimension)] = "$_id." + getMemberName(dimension)
	}

	for _, measure := range measures {
		projectStage["$project"].(bson.M)[getMemberName(measure)] = 1
	}
	return projectStage
}
