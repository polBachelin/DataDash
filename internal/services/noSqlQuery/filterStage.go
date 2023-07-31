package noSqlQuery

import (
	"dashboard/internal/services/query"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type FilterTypeFunc func(member string, value string) bson.M

var FilterTypes = map[string]FilterTypeFunc{
	"equals": FilterEquals,
}

func FilterEquals(member string, value string) bson.M {
	v := false
	if value == "No" || value == "false" {
		v = false
	} else if value == "Yes" || value == "true" {
		v = true
	}
	return bson.M{member: v}
}

func BuildFilterStage(filter query.Filter) (bson.M, error) {
	member := strings.Split(filter.Member, ".")[1]

	filterFunc, ok := FilterTypes[filter.Operator]
	if !ok {
		return bson.M{}, errors.New("no filter with type :" + filter.Operator)
	}
	stage := filterFunc(member, filter.Values[0])
	return bson.M{"$match": stage}, nil
}

func BuildAllFilters(filters []query.Filter) ([]bson.M, error) {
	filterStages := make([]bson.M, len(filters))

	for i, filter := range filters {
		filterStage, err := BuildFilterStage(filter)
		if err != nil {
			return []bson.M{}, err
		}
		filterStages[i] = filterStage
	}
	return filterStages, nil
}
