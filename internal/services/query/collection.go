package query

import blockService "dashboard/internal/services/block"

func FindCollectionName(dimensions []string, join *blockService.Join) string {

	if join != nil {
		for _, dimension := range dimensions {
			collectionName := getBlockName(dimension)
			if join.Name == collectionName {
				return collectionName
			}
		}
	}
	return getBlockName(dimensions[0])
}
