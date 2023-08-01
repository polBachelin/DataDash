package query

import (
	"dashboard/internal/database"
	"dashboard/internal/services/block"
	"dashboard/internal/services/sqlStages"
	"strings"
)

type Query struct {
	Measures       []string        `json:"measures"`
	Dimensions     []string        `json:"dimensions"`
	Filters        []Filter        `json:"filters"`
	TimeDimensions []TimeDimension `json:"time_dimensions"`
	Limit          int             `json:"limit"`
	Offset         int             `json:"offset"`
	Order          Order           `json:"order"`
}

type Filter struct {
	Member   string   `json:"member"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

type TimeDimension struct {
	Dimension   string   `json:"dimension"`
	DateRange   []string `json:"date_range"`
	Granularity string   `json:"granularity"`
}

type Order struct {
	DimensionName  []string `json:"dimension_name"`
	DimensionOrder []string `json:"dimension_order"`
	MeasureName    []string `json:"measure_name"`
	MeasureOrder   []string `json:"measure_order"`
}

type QueryService struct {
	Query Query
	Db    database.IDatabase
}

func NewQueryService(q Query, db database.IDatabase) *QueryService {
	return &QueryService{Query: q, Db: db}
}

func AddSelectToString(members []string, f func(string, *block.BlockData) string, res *strings.Builder) {
	memberLen := len(members)
	for i, m := range members {
		blockData := block.GetBlockFromName(GetBlockName(m))
		s := sqlStages.GenerateMeasureSelect(m, blockData)
		res.WriteString(s)
		if i+1 < memberLen {
			res.WriteRune(',')
		}
	}
}

func (service *QueryService) GenerateSelectStage() string {
	var result strings.Builder

	result.WriteString("SELECT ")
	AddSelectToString(service.Query.Measures, sqlStages.GenerateMeasureSelect, &result)
	AddSelectToString(service.Query.Dimensions, sqlStages.GenerateDimensionSelect, &result)
	return result.String()
}

func (service *QueryService) ParseQuery() (string, error) {
	//base := "SELECT %s"
	return "", nil
}
