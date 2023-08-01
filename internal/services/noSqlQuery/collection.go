package noSqlQuery

import (
	blockService "dashboard/internal/services/block"
	"dashboard/internal/services/query"
	"log"

	"golang.org/x/exp/slices"
)

func FindCollectionName(dimensions []string, join *blockService.Join) string {

	if join != nil {
		for _, dimension := range dimensions {
			collectionName := query.GetBlockName(dimension)
			if join.Name == collectionName {
				log.Println("Collection name: ", collectionName)
				index := slices.IndexFunc(dimensions, func(data string) bool {
					return data != dimension
				})
				return query.GetBlockName(dimensions[index])
			}
		}
	}
	log.Println("Collection name: ", query.GetBlockName(dimensions[0]))
	return query.GetBlockName(dimensions[0])
}
