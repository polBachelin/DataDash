package noSqlQuery

import (
	"dashboard/internal/services/block"
	blockService "dashboard/internal/services/block"
	"log"

	"golang.org/x/exp/slices"
)

func FindCollectionName(dimensions []string, join *blockService.Join) string {

	if join != nil {
		for _, dimension := range dimensions {
			collectionName := block.GetBlockName(dimension)
			if join.Name == collectionName {
				log.Println("Collection name: ", collectionName)
				index := slices.IndexFunc(dimensions, func(data string) bool {
					return data != dimension
				})
				return block.GetBlockName(dimensions[index])
			}
		}
	}
	log.Println("Collection name: ", block.GetBlockName(dimensions[0]))
	return block.GetBlockName(dimensions[0])
}
