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

func (query *Query) GenerateSelectStage() string {
	var result strings.Builder

	result.WriteString("SELECT ")
	AddSelectToString(query.Dimensions, sqlStages.GenerateDimensionSelect, &result)
	AddSelectToString(query.Measures, sqlStages.GenerateMeasureSelect, &result)
	return result.String()
}

func (query *Query) GetParentTableName() string {
	if len(query.Measures) > 0 {
		return block.GetBlockFromName(GetBlockName(query.Measures[0])).Table
	}
	if len(query.Dimensions) > 0 {
		return block.GetBlockFromName(GetBlockName(query.Dimensions[0])).Table
	}
	return ""
}

func GetBlockThatHasJoin(name string) *block.BlockData {
	blockInstance := block.GetInstance().Blocks

	for _, fileData := range blockInstance {
		for _, block := range fileData.Blocks {
			for _, blockJoin := range block.Joins {
				if blockJoin.Name == name {
					return &block
				}
			}
		}
	}
	return nil
}

func (query *Query) GenerateLeftJoinStage() string {
	blockInstance := block.GetInstance().Blocks

	for i, measure := range query.Measures {
		dimension := query.Dimensions[i]
		if strings.HasPrefix(measure, GetBlockName(dimension)) {
			continue
		}
		measureBlock := block.GetBlockFromName(GetBlockName(measure))
		dimensionBlock := block.GetBlockFromName(GetBlockName(dimension))
		if len(measureBlock.Joins) == 0 || len(dimensionBlock.Joins) == 0 {
			return ""
		}

	}
}

func (query *Query) GenerateFromStage() string {
	var result strings.Builder

	result.WriteString("FROM ")
	result.WriteString(query.GetParentTableName())
	if HasTwoDifferentBlocks(query.Dimensions, query.Measures) {
		result.WriteString(query.GenerateLeftJoinStage())
	}
	return result.String()
}

func (service *QueryService) ParseQuery() (string, error) {
	//base := "SELECT %s"
	selectStage := service.Query.GenerateSelectStage()

	return "", nil
}
