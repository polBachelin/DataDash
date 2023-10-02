package query

import (
	"dashboard/internal/database"
	"dashboard/internal/services/block"
	"dashboard/internal/services/sqlStages"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type Query struct {
	Measures       []string        `json:"measures"`
	Dimensions     []string        `json:"dimensions"`
	Filters        []Filter        `json:"filters"`
	TimeDimensions []TimeDimension `json:"timeDimensions"`
	Limit          int             `json:"limit"`
	Offset         int             `json:"offset"`
	Order          [][]string      `json:"order"`
}

type Filter struct {
	Member   string   `json:"member"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

type TimeDimension struct {
	Dimension   string   `json:"dimension"`
	DateRange   []string `json:"dateRange"`
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

func AddSelectToString(member string, genFunc func(string, *block.BlockData) (string, error)) (string, error) {
	blockData := block.GetBlockFromName(block.GetBlockName(member))
	if blockData == nil {
		return "", fmt.Errorf("no block with name %v", block.GetBlockName(member))
	}
	selectString, err := genFunc(GetMemberName(member), blockData)
	if err != nil {
		return "", err
	}
	return selectString, nil
}

func (query *Query) GenerateSelectStage() ([]string, error) {
	var result []string

	for _, measure := range query.Measures {
		selectString, err := AddSelectToString(measure, sqlStages.GenerateMeasureSql)
		if err != nil {
			return nil, err
		}
		result = append(result, selectString)
	}
	for _, dimension := range query.Dimensions {
		selectString, err := AddSelectToString(dimension, sqlStages.GenerateDimensionSelect)
		if err != nil {
			return nil, err
		}
		result = append(result, selectString)

	}
	return result, nil
}

func (query *Query) GetStartAndTargetTables() (*block.BlockData, []string) {
	if len(query.Measures) > 0 {
		b := block.GetBlockFromName(block.GetBlockName(query.Measures[0]))
		var t []string
		if len(query.Dimensions) == 0 {
			t = block.GetAllBlockNamesDifferent(b.Name, query.Measures)
		} else {
			t = block.GetAllBlockNamesDifferent(b.Name, query.Dimensions)
		}
		return b, t
	}
	if len(query.Dimensions) > 0 {
		b := block.GetBlockFromName(block.GetBlockName(query.Dimensions[0]))
		var t []string
		if len(query.Measures) == 0 {
			t = block.GetAllBlockNamesDifferent(b.Name, query.Dimensions)
		} else {
			t = block.GetAllBlockNamesDifferent(b.Name, query.Measures)
		}
		return b, t
	}
	return nil, nil
}

func (query *Query) GenerateLeftJoinStage(graph *block.JoinGraph) []string {
	startTableName, targetTableNames := query.GetStartAndTargetTables()
	return sqlStages.GetLeftJoinPath(startTableName, targetTableNames, graph)
}

func (query *Query) GenerateFromStage(graph *block.JoinGraph) string {
	var result strings.Builder

	result.WriteString(" FROM ")
	var path []string
	if HasTwoDifferentBlocks(query.Dimensions, query.Measures) {
		path = query.GenerateLeftJoinStage(graph)
		if len(path) == 0 {
			return ""
		}
		firstBlock := block.GetBlockFromName(path[0])
		result.WriteString(fmt.Sprintf("%v as %v", firstBlock.Table, firstBlock.Name))
		for i := 1; i < len(path); i++ {
			fromVertex := graph.Vertices[path[i-1]]
			toVertex := graph.Vertices[path[i]]

			joinParent, err := block.GetBlockJoinFromName(toVertex.Val.Name, fromVertex.Val)
			if err != nil {
				joinParent, _ = block.GetBlockJoinFromName(fromVertex.Val.Name, toVertex.Val)
			}
			result.WriteString(fmt.Sprintf(" LEFT JOIN %s as %s ON %s.%s = %s.%s",
				toVertex.Val.Table, toVertex.Val.Name, toVertex.Val.Name,
				joinParent.LocalField, fromVertex.Val.Name, joinParent.ForeignField))
		}
	} else {
		parentTable, _ := query.GetStartAndTargetTables()
		result.WriteString(fmt.Sprintf("%s as %s", parentTable.Table, parentTable.Name))
	}
	return result.String()
}

func (query *Query) CountsInMeasures() (int, error) {
	n := 0
	for _, measure := range query.Measures {
		b := block.GetBlockFromName(block.GetBlockName(measure))
		m, err := block.GetMeasureFromBlock(b, GetMemberName(measure))
		if err != nil {
			return 0, err
		}
		if MeasureIsAggregated(m) {
			n++
		}
	}
	return n, nil
}

func (query *Query) FilterHasAggregated() bool {
	for _, f := range query.Filters {
		b := block.GetBlockFromName(block.GetBlockName(f.Member))
		m, err := block.GetMeasureFromBlock(b, GetMemberName(f.Member))
		if err != nil {
			return false
		}
		if MeasureIsAggregated(m) {
			return true
		}
	}
	return false
}

func (query *Query) GenerateGroupByStage(totalSelect int) string {
	var result strings.Builder

	result.WriteString(" GROUP BY ")
	measureLen, _ := query.CountsInMeasures()
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

func (query *Query) GenerateOrderStage() ([]string, error) {
	var i int
	var result []string

	for _, order := range query.Order {
		if len(order) < 2 {
			return nil, fmt.Errorf("order needs to contain two values [memberName, order]")
		}
		if i = slices.Index(query.Measures, order[0]); i == -1 {
			i = slices.Index(query.Dimensions, order[0]) + len(query.Measures)
		}
		if i == -1 {
			return nil, fmt.Errorf("order does not contain a member present in the query")
		}
		if !strings.EqualFold(order[1], "asc") && !strings.EqualFold(order[1], "desc") {
			return nil, fmt.Errorf("order is not asc or desc")
		}
		result = append(result, fmt.Sprintf("%v %v", i+1, strings.ToUpper(order[1])))
	}
	return result, nil
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
	if timeD.Dimension == "" {
		return "", "", fmt.Errorf("no dimension in the timeDimension")
	}
	b := block.GetBlockFromName(block.GetBlockName(timeD.Dimension))
	memberName := GetMemberName(timeD.Dimension)
	dimension, _ := block.GetDimensionFromBlock(b, memberName)
	if dimension == nil {
		return "", "", fmt.Errorf("dimension not found in block %v", b.Name)
	}
	return fmt.Sprintf("date_trunc('%v', (%v.%v :: timestamptz AT TIME ZONE 'UTC')) \"%v_%v_%v\"", timeD.Granularity, b.Name, dimension.Sql, b.Name, memberName, timeD.Granularity), fmt.Sprintf("(%v.%v >= '%v' :: timestamptz AND %v.%v <= '%v' :: timestamptz)", b.Name, dimension.Sql, timeD.DateRange[0], b.Name, dimension.Sql, timeD.DateRange[1]), nil
}

type FilterContext struct {
	isMember bool
	Sql      string
}

func (query *Query) GenerateFilterMap() (map[string]FilterContext, error) {
	filterMap := make(map[string]FilterContext)
	for _, filter := range query.Filters {
		b := block.GetBlockFromName(block.GetBlockName(filter.Member))
		f, isHaving, err := sqlStages.GenerateFilter(b, filter.Values, GetMemberName(filter.Member), filter.Operator)
		if err != nil {
			return nil, err
		}
		filterMap[filter.Member] = FilterContext{isMember: isHaving, Sql: f}
	}
	return filterMap, nil
}
