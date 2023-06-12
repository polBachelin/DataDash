package query

import "go.mongodb.org/mongo-driver/bson"

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
