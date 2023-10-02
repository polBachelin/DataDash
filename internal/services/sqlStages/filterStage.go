package sqlStages

import (
	"dashboard/internal/services/block"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type FilterTypeFunc func(blockData *block.BlockData, member string, values []string, measure bool) (string, error)

var FilterTypes = map[string]FilterTypeFunc{
	"equals":     FilterEquals,
	"gt":         FilterGreater,
	"gte":        FilterGreaterEquals,
	"lt":         FilterLess,
	"lte":        FilterLessEquals,
	"dateEquals": FilterDateEquals,
	"beforeDate": FilterBeforeDate,
}

func FilterMathOperation(blockData *block.BlockData, member string, values []string, measure bool, operation string) (string, error) {
	var result strings.Builder

	if measure {
		sql, err := GenerateMeasureSql(member, blockData)
		if err != nil {
			return "", fmt.Errorf("Could not generate measure sql: %v", err)
		}
		result.WriteString(sql)
		result.WriteString(fmt.Sprintf("%s %v", operation, values[0]))
	} else {
		result.WriteString(fmt.Sprintf("%v.%v %s '%v'", blockData.Name, member, operation, values[0]))
	}
	return result.String(), nil
}

func FilterDateEquals(blockData *block.BlockData, member string, values []string, measure bool) (string, error) {
	return fmt.Sprintf("cast(%v.%v as date)  = cast(%v as date)", blockData.Name, member, values[0]), nil
}

func FilterBeforeDate(blockData *block.BlockData, member string, values []string, measure bool) (string, error) {
	return fmt.Sprintf("cast(%v.%v as date)  < cast(%v as date)", blockData.Name, member, values[0]), nil
}

func FilterEquals(blockData *block.BlockData, member string, values []string, measure bool) (string, error) {
	// TODO: ADD VALUES AS ($1, $2....) to be able to save the query to cache this will need to change how the query is sent to db
	return FilterMathOperation(blockData, member, values, measure, "=")
}

func FilterGreater(blockData *block.BlockData, member string, values []string, measure bool) (string, error) {
	return FilterMathOperation(blockData, member, values, measure, ">")
}

func FilterGreaterEquals(blockData *block.BlockData, member string, values []string, measure bool) (string, error) {
	return FilterMathOperation(blockData, member, values, measure, ">=")
}

func FilterLess(blockData *block.BlockData, member string, values []string, measure bool) (string, error) {
	return FilterMathOperation(blockData, member, values, measure, "<")

}

func FilterLessEquals(blockData *block.BlockData, member string, values []string, measure bool) (string, error) {
	return FilterMathOperation(blockData, member, values, measure, "<=")

}

func GenerateFilter(memberBlock *block.BlockData, values []string, member, operator string) (string, bool, error) {
	if _, ok := FilterTypes[operator]; !ok {
		return "", false, fmt.Errorf("this operator does not exist %s", operator)
	}
	if slices.ContainsFunc(memberBlock.Measures, func(data block.Measures) bool { return data.Name == member }) {
		res, err := FilterTypes[operator](memberBlock, member, values, true)
		if err != nil {
			return "", true, err
		}
		return res, true, nil //TODO: improvement could be made here, I don't like the fact of just puttin a boolean like that
	}
	res, err := FilterTypes[operator](memberBlock, member, values, false)
	if err != nil {
		return "", false, err
	}
	return res, false, nil
}
