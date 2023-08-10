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
		blockData := block.GetBlockFromName(GetBlockName(m))
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
	if len(query.Dimensions) > 0 {
		result.WriteRune(',')
	}
	AddSelectToString(query.Dimensions, sqlStages.GenerateDimensionSelect, &result)
	return result.String()
}

func (query *Query) GetParentTableName() string {
	if len(query.Measures) > 0 {
		b := block.GetBlockFromName(GetBlockName(query.Measures[0]))
		return b.Name
	}
	if len(query.Dimensions) > 0 {
		b := block.GetBlockFromName(GetBlockName(query.Dimensions[0]))
		return b.Name
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
	targetTableName := "Status_name" //Need to get this from query

	log.Printf("startTableName=%s", startTableName)
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

	log.Printf("Join clause path = %s", path)
	for i := len(path) - 1; i >= 1; i-- {
		fromVertex := graph.Vertices[path[i]]
		toVertex := graph.Vertices[path[i-1]]

		joinParent, err := block.GetBlockJoinFromName(toVertex.Val.Name, *fromVertex.Val)
		if err != nil {
			joinParent, err = block.GetBlockJoinFromName(fromVertex.Val.Name, *toVertex.Val)
		}
		joins.WriteString(fmt.Sprintf(" LEFT JOIN %s as %s ON %s.%s = %s.%s", toVertex.Val.Table, toVertex.Val.Name, toVertex.Val.Name, joinParent.LocalField, joinParent.Name, joinParent.ForeignField))
	}
	return joins.String()
}

func (query *Query) GenerateFromStage(graph *block.JoinGraph) string {
	var result strings.Builder

	result.WriteString("FROM ")
	result.WriteString(query.GetParentTableName())
	if HasTwoDifferentBlocks(query.Dimensions, query.Measures) {
		result.WriteString(query.GenerateLeftJoinStage(graph))
	}
	return result.String()
}

func (query *Query) GenerateGroupByStage() string {
	var result strings.Builder

	result.WriteString("GROUP BY ")
	n := len(query.Measures) + 1
	for i := range query.Dimensions {
		result.WriteString(fmt.Sprintf("%d", i+n))
		if i < len(query.Dimensions)-1 {
			result.WriteRune(',')
		}
	}
	return result.String()
}

func (service *QueryService) ParseQuery() (string, error) {
	selectStage := service.Query.GenerateSelectStage()
	fromStage := service.Query.GenerateFromStage(service.JoinGraph)
	groupByStage := service.Query.GenerateGroupByStage()

	log.Println("GENERATED SQL : ", selectStage, fromStage, groupByStage)
	return "", nil
}
