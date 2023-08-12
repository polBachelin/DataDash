package query

import (
	"dashboard/internal/database"
	"dashboard/internal/services/block"
	"dashboard/internal/services/sqlStages"
	"fmt"
	"log"
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
	JoinGraph *block.JoinGraph
	Db        database.IDatabase
}

func NewQueryService(q Query, db database.IDatabase, joinGraph *block.JoinGraph) *QueryService {
	return &QueryService{Query: q, Db: db, JoinGraph: joinGraph}
}

func AddSelectToString(members []string, genFunc func(string, *block.BlockData) string, res *strings.Builder) {
	memberLen := len(members)
	for i, m := range members {
		blockData := block.GetBlockFromName(block.GetBlockName(m))
		s := genFunc(GetMemberName(m), blockData)
		res.WriteString(s)
		log.Println(s)
		if i+1 < memberLen {
			res.WriteRune(',')
		}
	}
}

func (query *Query) GenerateSelectStage() string {
	var result strings.Builder

	result.WriteString("SELECT ")
	AddSelectToString(query.Measures, sqlStages.GenerateMeasureSelect, &result)
	if len(query.Dimensions) > 0 && len(query.Measures) > 0 {
		result.WriteRune(',')
	}
	AddSelectToString(query.Dimensions, sqlStages.GenerateDimensionSelect, &result)
	return result.String()
}

func (query *Query) GetStartAndTargetTables() (string, []string) {
	if len(query.Measures) > 0 {
		b := block.GetBlockFromName(block.GetBlockName(query.Measures[0]))
		t := block.GetAllBlockNamesDifferent(b.Name, query.Dimensions)
		return b.Name, t
	}
	if len(query.Dimensions) > 0 {
		b := block.GetBlockFromName(block.GetBlockName(query.Dimensions[0]))
		t := block.GetAllBlockNamesDifferent(b.Name, query.Measures)
		return b.Name, t
	}
	return "", nil
}

func (query *Query) GenerateLeftJoinStage(graph *block.JoinGraph) string {
	startTableName, targetTableNames := query.GetStartAndTargetTables()

	for _, targetTable := range targetTableNames {
		if startVertex, found := graph.Vertices[startTableName]; found {
			path, relationshipFound := graph.FindJoinPath(startVertex, targetTable)
			if relationshipFound {
				joins := query.GenerateJoinClause(path, graph)
				return joins
			}
		}
	}
	return ""
}

func (query *Query) GenerateJoinClause(path []string, graph *block.JoinGraph) string {
	var joins strings.Builder

	for i := len(path) - 1; i >= 1; i-- {
		fromVertex := graph.Vertices[path[i]]
		toVertex := graph.Vertices[path[i-1]]

		joinParent, err := block.GetBlockJoinFromName(toVertex.Val.Name, fromVertex.Val)
		if err != nil {
			joinParent, _ = block.GetBlockJoinFromName(fromVertex.Val.Name, toVertex.Val)
		}
		joins.WriteString(fmt.Sprintf(" LEFT JOIN %s as %s ON %s.%s = %s.%s", toVertex.Val.Table, toVertex.Val.Name, toVertex.Val.Name, joinParent.LocalField, fromVertex.Val.Name, joinParent.ForeignField))
	}
	return joins.String()
}

func (query *Query) GenerateFromStage(graph *block.JoinGraph) string {
	var result strings.Builder

	result.WriteString(" FROM ")
	parentTable, _ := query.GetStartAndTargetTables()
	result.WriteString(parentTable)
	if HasTwoDifferentBlocks(query.Dimensions, query.Measures) {
		result.WriteString(query.GenerateLeftJoinStage(graph))
	}
	return result.String()
}

func (query *Query) GenerateGroupByStage() string {
	var result strings.Builder

	result.WriteString(" GROUP BY ")
	n := len(query.Measures) + 1
	for i := range query.Dimensions {
		result.WriteString(fmt.Sprintf("%d", i+n))
		if i < len(query.Dimensions)-1 {
			result.WriteRune(',')
		}
	}
	return result.String()
}

func (query *Query) GenerateFilterStage() string {
	var result strings.Builder

	for _, filter := range query.Filters {
		b := block.GetBlockFromName(block.GetBlockName(filter.Operator))
		f := sqlStages.GenerateFilters(b, filter.Values, filter.Operator)
		result.WriteString(f)
	}
	return result.String()
}

func (service *QueryService) ParseQuery() ([]map[string]interface{}, error) {
	var sqlQuery strings.Builder

	sqlQuery.WriteString(service.Query.GenerateSelectStage())
	sqlQuery.WriteString(service.Query.GenerateFromStage(service.JoinGraph))
	sqlQuery.WriteString(service.Query.GenerateGroupByStage())
	sqlQuery.WriteString(service.Query.GenerateFilterStage())
	sqlResult, err := service.Db.ExecuteQuery(sqlQuery.String())
	if err != nil {
		return nil, err
	}
	resJson, err := service.Db.QueryResultToJson(sqlResult)
	if err != nil {
		return nil, err
	}
	return resJson, nil
}
