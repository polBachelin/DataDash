package sqlStages

import "dashboard/internal/services/block"

type FilterTypeFunc func(member string, values []string) string

var FilterTypes = map[string]FilterTypeFunc{
	"equals": FilterEquals,
}

func FilterEquals(member string, values []string) string {
	return ""
}

func GenerateFilters(memberBlock *block.BlockData, values []string, operator string) string {

	return ""
}
