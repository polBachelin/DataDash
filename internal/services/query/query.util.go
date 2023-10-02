package query

import (
	"dashboard/internal/services/block"
	"strings"
)

func FindBlockWithJoin(dimensions []string) *block.Join {
	for i, dimension := range dimensions {
		block := block.GetBlockFromName(block.GetBlockName(dimension))
		for _, join := range block.Joins {
			if HasBlockName(dimensions[i+1:], join.Name) {
				return &join
			}
		}
	}
	return nil
}

func HasBlockName(dimensions []string, targetDimension string) bool {
	for _, dimension := range dimensions {
		if block.GetBlockName(dimension) == targetDimension {
			return true
		}
	}
	return false
}

func GetMemberName(dimension string) string {
	parts := strings.Split(dimension, ".")
	return parts[1]
}

func HasTwoDifferentBlocks(dimensions []string, measures []string) bool {
	blockCount := make(map[string]int)

	for _, dimension := range dimensions {
		blockName := block.GetBlockName(dimension)
		blockCount[blockName]++
	}
	for _, measure := range measures {
		blockName := block.GetBlockName(measure)
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

func GetMeasureType(measure string) string {
	memberName := GetMemberName(measure)
	b := block.GetBlockFromName(block.GetBlockName(measure))
	for _, m := range b.Measures {
		if m.Name == memberName {
			return m.Type
		}
	}
	return ""
}

func GetDimensionType(dimension string) string {
	memberName := GetMemberName(dimension)
	b := block.GetBlockFromName(block.GetBlockName(dimension))
	for _, m := range b.Dimensions {
		if m.Name == memberName {
			return m.Type
		}
	}
	return ""
}
<<<<<<< HEAD
=======

func StringIsAggregateFunction(sql string) bool {
	functions := []string{"count", "sum", "min", "max", "avg"}

	sql = strings.ToLower(sql)
	for _, f := range functions {
		if strings.HasPrefix(sql, f) {
			return true
		}
	}
	return false
}

func GetTitle(memberName string) string {
	var res strings.Builder

	parts := strings.Split(memberName, ".")
	res.WriteString(parts[0])
	res.WriteRune(' ')
	res.WriteString(strings.Title(parts[1]))
	return res.String()
}

func GetShortTitle(memberName string) string {
	m := GetMemberName(memberName)
	return strings.Title(m)
}

func MeasureIsAggregated(m *block.Measures) bool {
	if m.Type == "count" || m.Type == "sum" || StringIsAggregateFunction(m.Sql) {
		return true
	}
	return false
}
>>>>>>> dev
