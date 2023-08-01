package query

import (
	"strings"
)

func GetBlockName(dimension string) string {
	parts := strings.Split(dimension, ".")
	return parts[0]
}

func HasBlockName(dimensions []string, targetDimension string) bool {
	for _, dimension := range dimensions {
		if GetBlockName(dimension) == targetDimension {
			return true
		}
	}
	return false
}

func GetMemberName(dimension string) string {
	parts := strings.Split(dimension, ".")
	return parts[1]
}

func HasTwoDifferentBlocks(dimensions []string) bool {
	blockCount := make(map[string]int)

	for _, dimension := range dimensions {
		blockName := GetBlockName(dimension)
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

func ContainsDimension(dimensions []string, dimension string) bool {
	for _, d := range dimensions {
		if d == dimension {
			return true
		}
	}
	return false
}
