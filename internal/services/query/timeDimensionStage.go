package query

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func BuildTimeDimension(timeDimension TimeDimension) bson.M {
	member := strings.Split(timeDimension.Dimension, ".")[1]
	return bson.M{"$match": bson.M{member: bson.M{"$gte": timeDimension.DateRange[0], "$lte": timeDimension.DateRange[1]}}}
}

// TODO check if I should enable the fact to add time dimensions for a join
// TODO not using granularity need to see how I can take advantage of that
func BuildAllTimeDimensions(timeDimensions []TimeDimension) ([]bson.M, error) {
	timeDimensionStages := make([]bson.M, len(timeDimensions))

	for i, d := range timeDimensions {
		stage := BuildTimeDimension(d)
		timeDimensionStages[i] = stage
	}
	return timeDimensionStages, nil
}

func getDateFormat(granularity string) string {
	switch granularity {
	case "day":
		return "%Y-%m-%d"
	case "month":
		return "%Y-%m"
	case "year":
		return "%Y"
	default:
		return "%Y-%m-%d"
	}
}
