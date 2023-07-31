package noSqlQuery

import (
	"strings"
)

func getBlockName(dimension string) string {
	parts := strings.Split(dimension, ".")
	return parts[0]
}

func hasBlockName(dimensions []string, targetDimension string) bool {
	for _, dimension := range dimensions {
		if getBlockName(dimension) == targetDimension {
			return true
		}
	}
	return false
}

func getMemberName(dimension string) string {
	parts := strings.Split(dimension, ".")
	return parts[1]
}

func hasTwoDifferentBlocks(dimensions []string) bool {
	blockCount := make(map[string]int)

	for _, dimension := range dimensions {
		blockName := getBlockName(dimension)
		blockCount[blockName]++
	}

	blockCountSize := 0
	for _, count := range blockCount {
		if count > 0 {
			blockCountSize++
		}
	}

	return blockCountSize >= 2
}

func containsDimension(dimensions []string, dimension string) bool {
	for _, d := range dimensions {
		if d == dimension {
			return true
		}
	}
	return false
}
