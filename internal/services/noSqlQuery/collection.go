package noSqlQuery

import (
	blockService "dashboard/internal/services/block"
	"log"

	"golang.org/x/exp/slices"
)

func FindCollectionName(dimensions []string, join *blockService.Join) string {

	if join != nil {
		for _, dimension := range dimensions {
			collectionName := getBlockName(dimension)
			if join.Name == collectionName {
				log.Println("Collection name: ", collectionName)
				index := slices.IndexFunc(dimensions, func(data string) bool {
					return data != dimension
				})
				return getBlockName(dimensions[index])
			}
		}
	}
	log.Println("Collection name: ", getBlockName(dimensions[0]))
	return getBlockName(dimensions[0])
}
