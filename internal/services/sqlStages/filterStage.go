package sqlStages

import (
	"dashboard/internal/services/block"
	"fmt"
	"log"
	"strings"

	"golang.org/x/exp/slices"
)

type FilterTypeFunc func(blockData *block.BlockData, member string, values []string, measure bool) string

var FilterTypes = map[string]FilterTypeFunc{
	"equals": FilterEquals,
	"gt":     FilterGreater,
	"gte":    FilterGreaterEquals,
	"lt":     FilterLess,
	"lte":    FilterLessEquals,
}

func FilterMathOperation(blockData *block.BlockData, member string, values []string, measure bool, operation string) string {
	var result strings.Builder

	if measure {
		sql := GenerateMeasureSql(member, blockData)
		result.WriteString(sql)
		result.WriteString(fmt.Sprintf("%s %v", operation, values[0]))
	} else {
		result.WriteString(fmt.Sprintf("%v.%v %s '%v'", blockData.Name, member, operation, values[0]))
		log.Println(result.String())
	}
	return result.String()
}

func FilterEquals(blockData *block.BlockData, member string, values []string, measure bool) string {
	// TODO: ADD VALUES AS ($1, $2....) to be able to save the query to cache this will need to change how the query is sent to db
	return FilterMathOperation(blockData, member, values, measure, "=")
}

func FilterGreater(blockData *block.BlockData, member string, values []string, measure bool) string {
	return FilterMathOperation(blockData, member, values, measure, ">")
}

func FilterGreaterEquals(blockData *block.BlockData, member string, values []string, measure bool) string {
	return FilterMathOperation(blockData, member, values, measure, ">=")
}

func FilterLess(blockData *block.BlockData, member string, values []string, measure bool) string {
	return FilterMathOperation(blockData, member, values, measure, "<")

}

func FilterLessEquals(blockData *block.BlockData, member string, values []string, measure bool) string {
	return FilterMathOperation(blockData, member, values, measure, "<=")

}

func GenerateFilter(memberBlock *block.BlockData, values []string, member, operator string) (string, bool, error) {
	if _, ok := FilterTypes[operator]; !ok {
		return "", false, fmt.Errorf("this operator does not exist %s", operator)
	}
	log.Println(memberBlock)
	if slices.ContainsFunc(memberBlock.Measures, func(data block.Measures) bool { return data.Name == member }) {
		return FilterTypes[operator](memberBlock, member, values, true), true, nil //TODO: improvement could be made here, I don't like the fact of just puttin a boolean like that
	}
	return FilterTypes[operator](memberBlock, member, values, false), false, nil
}
