package noSqlQuery

import (
	"dashboard/internal/services/query"

	"go.mongodb.org/mongo-driver/bson"
)

func generateOrderStage(order [][]string) bson.M {
	sort := bson.M{}
	for _, member := range order {
		sort[query.GetMemberName(member[0])] = getOrderType(member[1])
	}
	return bson.M{"$sort": sort}
}

func getOrderType(orderType string) int {
	switch orderType {
	case "asc":
		return 1
	case "desc":
		return -1
	default:
		return 1
	}
}
