package query

import (
	"dashboard/internal/database"
	"dashboard/internal/services/block"
	"dashboard/internal/services/sqlStages"
	"fmt"
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
	Query     Query
	JoinGraph block.JoinGraph
	Db        database.IDatabase
}

func NewQueryService(q Query, db database.IDatabase, joinGraph block.JoinGraph) *QueryService {
	return &QueryService{Query: q, Db: db, JoinGraph: joinGraph}
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

func (query *Query) GenerateLeftJoinStage(graph *block.JoinGraph) string {
	startTableName := query.GetParentTableName()
	targetTableName := "Status_name"

	if startVertex, found := graph.Vertices[startTableName]; found {
		path, relationshipFound := graph.FindJoinPath(startVertex, targetTableName)
		if relationshipFound {
			joins := query.GenerateJoinClause(path, graph)
			return joins
		}
	}
	return ""
}

func (query *Query) GenerateJoinClause(path []string, graph *block.JoinGraph) string {
	var joins strings.Builder

	for i := len(path) - 1; i >= 0; i-- {
		fromVertex := graph.Vertices[path[i]]
		toVertex := graph.Vertices[path[i-1]]

		joinParent, err := block.GetBlockJoinFromName(toVertex.Val.Name, *fromVertex.Val)
		if err != nil {
			joinParent, err = block.GetBlockJoinFromName(fromVertex.Val.Name, *toVertex.Val) //Order_status
		}
		joins.WriteString(fmt.Sprintf("LEFT JOIN %s as %s ON %s.%s = %s.%s", toVertex.Val.Table, toVertex.Val.Name, toVertex.Val.Name, joinParent.LocalField, joinParent.Name, joinParent.ForeignField))
	}
	return joins.String()
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
	selectStage := service.Query.GenerateSelectStage()

	return "", nil
}
