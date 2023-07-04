package query

import (
	blockService "dashboard/internal/services/block"

	"go.mongodb.org/mongo-driver/bson"
)

func BuildLookupStage(join blockService.Join) bson.M {
	return bson.M{
		"$lookup": bson.M{
			"from":         join.Name,
			"localField":   join.LocalField,
			"foreignField": join.ForeignField,
			"as":           join.Name,
		}}
}
