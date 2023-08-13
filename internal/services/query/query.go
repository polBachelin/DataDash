package query

import (
	"dashboard/internal/database"
	"dashboard/internal/services/block"
	"dashboard/internal/services/sqlStages"
	"fmt"
	"log"
	"strings"

	"golang.org/x/exp/slices"
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

var MeasureFilters = []string{"equals", "notEquals", "gte", "gt", "lt", "lte", "set", "notSet"}

func NewQueryService(q Query, db database.IDatabase, joinGraph *block.JoinGraph) *QueryService {
	return &QueryService{Query: q, Db: db, JoinGraph: joinGraph}
}

func AddSelectToString(member string, genFunc func(string, *block.BlockData) string) string {
	blockData := block.GetBlockFromName(block.GetBlockName(member))
	return genFunc(GetMemberName(member), blockData)
}

func (query *Query) GenerateSelectStage() []string {
	var result []string

	for _, measure := range query.Measures {
		result = append(result, AddSelectToString(measure, sqlStages.GenerateMeasureSql))
	}
	for _, dimension := range query.Dimensions {
		result = append(result, AddSelectToString(dimension, sqlStages.GenerateDimensionSelect))
	}
	return result
}

func (query *Query) GetStartAndTargetTables() (*block.BlockData, []string) {
	if len(query.Measures) > 0 {
		b := block.GetBlockFromName(block.GetBlockName(query.Measures[0]))
		t := block.GetAllBlockNamesDifferent(b.Name, query.Dimensions)
		return b, t
	}
	if len(query.Dimensions) > 0 {
		b := block.GetBlockFromName(block.GetBlockName(query.Dimensions[0]))
		t := block.GetAllBlockNamesDifferent(b.Name, query.Measures)
		return b, t
	}
	return nil, nil
}

func (query *Query) GenerateLeftJoinStage(graph *block.JoinGraph) string {
	startTableName, targetTableNames := query.GetStartAndTargetTables()

	for _, targetTable := range targetTableNames {
		if startVertex, found := graph.Vertices[startTableName.Name]; found {
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
	result.WriteString(fmt.Sprintf("%v as %v", parentTable.Table, parentTable.Name))
	if HasTwoDifferentBlocks(query.Dimensions, query.Measures) {
		result.WriteString(query.GenerateLeftJoinStage(graph))
	}
	return result.String()
}

func (query *Query) GenerateGroupByStage(totalSelect int) string {
	var result strings.Builder

	result.WriteString(" GROUP BY ")
	measureLen := len(query.Measures)
	i := measureLen
	if measureLen == 0 {
		i = 1
	}
	for i <= totalSelect-measureLen {
		result.WriteString(fmt.Sprintf("%d", i+measureLen))
		if i <= totalSelect-measureLen-1 {
			result.WriteRune(',')
		}
		i++
	}
	return result.String()
}

func (query *Query) FilterHasMeasure() (bool, error) {
	for _, filter := range query.Filters {
		b := block.GetBlockFromName(block.GetBlockName(filter.Member))
		if slices.ContainsFunc(b.Measures, func(data block.Measures) bool { return data.Name == GetMemberName(filter.Member) }) {
			if !slices.Contains(MeasureFilters, filter.Operator) {
				return false, fmt.Errorf("operator %s cannot be applied to measures", filter.Operator)
			}
			return true, nil
		}
	}
	return false, nil
}

func (query *Query) GenerateLimitStage() string {
	return fmt.Sprintf(" LIMIT %v", query.Limit)
}

func (query *Query) GenerateOffsetStage() string {
	return fmt.Sprintf(" OFFSET %v", query.Offset)
}

func (query *Query) GenerateTimeDimensionStage(index int) (string, string, error) {
	timeD := query.TimeDimensions[index]
	if len(timeD.DateRange) < 2 {
		return "", "", fmt.Errorf("not enough dates in daterange")
	}
	b := block.GetBlockFromName(block.GetBlockName(timeD.Dimension))
	memberName := GetMemberName(timeD.Dimension)
	dimension, _ := block.GetDimensionFromBlock(b, memberName)
	return fmt.Sprintf("date_trunc('%v', (%v.%v :: timestamptz AT TIME ZONE 'UTC')) \"%v_%v_%v\"", timeD.Granularity, b.Name, dimension.Sql, b.Name, memberName, timeD.Granularity), fmt.Sprintf("(%v.%v >= '%v' :: timestamptz AND %v.%v <= '%v' :: timestamptz)", b.Name, dimension.Sql, timeD.DateRange[0], b.Name, dimension.Sql, timeD.DateRange[1]), nil
}

type FilterContext struct {
	isMember bool
	Sql      string
}

func (service *QueryService) BuildStage(stage []string, start string, seperator string) string {
	var result strings.Builder

	stageLen := len(stage)
	result.WriteString(start)
	for i, stage := range stage {
		result.WriteString(stage)
		if i < stageLen-1 {
			result.WriteString(seperator)
		}
	}
	return result.String()
}

func (query *Query) GenerateFilterMap() map[string]FilterContext {
	filterMap := make(map[string]FilterContext)
	for _, filter := range query.Filters {
		b := block.GetBlockFromName(block.GetBlockName(filter.Member))
		f, isHaving, _ := sqlStages.GenerateFilter(b, filter.Values, GetMemberName(filter.Member), filter.Operator)
		filterMap[filter.Member] = FilterContext{isMember: isHaving, Sql: f}
	}
	return filterMap
}

func (service *QueryService) FilterMapToArray(filtersMap map[string]FilterContext) []string {
	var result []string

	for _, value := range filtersMap {
		result = append(result, value.Sql)
	}
	return result
}

func (service *QueryService) ParseQuery() ([]map[string]interface{}, error) {
	var sqlQuery strings.Builder
	var whereStage []string

	selectStage := service.Query.GenerateSelectStage()
	filterMap := make(map[string]FilterContext)
	if len(service.Query.Filters) > 0 {
		filterMap = service.Query.GenerateFilterMap()
		for key, value := range filterMap {
			if !value.isMember {
				whereStage = append(whereStage, value.Sql)
				delete(filterMap, key)
			}
		}
	}
	if len(service.Query.TimeDimensions) > 0 {
		selectTimeD, whereTimeD, err := service.Query.GenerateTimeDimensionStage(0)
		if err != nil {
			return nil, err
		}
		selectStage = append(selectStage, selectTimeD)
		whereStage = append(whereStage, whereTimeD)
	}
	sqlQuery.WriteString(service.BuildStage(selectStage, "SELECT ", ", "))
	sqlQuery.WriteString(service.Query.GenerateFromStage(service.JoinGraph))
	sqlQuery.WriteString(service.BuildStage(whereStage, " WHERE ", " AND "))
	sqlQuery.WriteString(service.Query.GenerateGroupByStage(len(selectStage)))
	if len(service.Query.Filters) > 0 {
		havingStage := service.FilterMapToArray(filterMap)
		sqlQuery.WriteString(service.BuildStage(havingStage, " HAVING ", " AND "))
	}
	sqlQuery.WriteString(service.Query.GenerateLimitStage())
	sqlQuery.WriteString(service.Query.GenerateOffsetStage())
	sqlResult, err := service.Db.ExecuteQuery(sqlQuery.String())
	log.Println(sqlQuery.String())
	if err != nil {
		return nil, err
	}
	resJson, err := service.Db.QueryResultToJson(sqlResult)
	if err != nil {
		return nil, err
	}
	return resJson, nil
}
