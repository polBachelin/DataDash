package query

import "strings"

type BlockQuery struct {
	Measure    []string
	Dimensions []string
	Name       string
}

func buildBlockQueryData(arr []string) []string {
	blockQueryData := make([]string, len(arr))
	for i, v := range arr {
		n := strings.Split(v, ".")
		blockQueryData[i] = n[1]
	}
	return blockQueryData
}

func buildBlockQuery(dimensions []string, measures []string, blockName string) BlockQuery {
	blockQuery := BlockQuery{Name: blockName, Measure: []string{}, Dimensions: []string{}}
	blockQuery.Measure = buildBlockQueryData(measures)
	blockQuery.Dimensions = buildBlockQueryData(dimensions)
	return blockQuery
}
