package query

import (
	"strings"
)

type Annotations struct {
	Measures      map[string]Annotation `json:"measures"`
	Dimensions    map[string]Annotation `json:"dimensions"`
	TimeDimension map[string]Annotation `json:"timeDimensions"`
}

type Annotation struct {
	Title      string `json:"title"`
	ShortTitle string `json:"shortTitle"`
	Type       string `json:"type"`
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

func (service *QueryService) FilterMapToArray(filtersMap map[string]FilterContext) []string {
	var result []string

	for _, value := range filtersMap {
		result = append(result, value.Sql)
	}
	return result
}

// TODO: need to cut this function it is too long
func (service *QueryService) ParseQuery() ([]map[string]interface{}, error) {
	var sqlQuery strings.Builder
	var whereStage []string

	selectStage, err := service.Query.GenerateSelectStage()
	if err != nil {
		return nil, err
	}
	filterMap := make(map[string]FilterContext)
	if len(service.Query.Filters) > 0 {
		filterMap, err = service.Query.GenerateFilterMap()
		if err != nil {
			return nil, err
		}
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
	if len(whereStage) >= 1 {
		sqlQuery.WriteString(service.BuildStage(whereStage, " WHERE ", " AND "))
	}
	if len(selectStage) > 1 {
		sqlQuery.WriteString(service.Query.GenerateGroupByStage(len(selectStage)))
	}
	if service.Query.FilterHasAggregated() {
		havingStage := service.FilterMapToArray(filterMap)
		sqlQuery.WriteString(service.BuildStage(havingStage, " HAVING ", " AND "))
	}
	if len(service.Query.Order) > 0 {
		orderStage, err := service.Query.GenerateOrderStage()
		if err != nil {
			return nil, err
		}
		sqlQuery.WriteString(service.BuildStage(orderStage, " ORDER BY ", ","))
	}
	sqlQuery.WriteString(service.Query.GenerateLimitStage())
	sqlQuery.WriteString(service.Query.GenerateOffsetStage())
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

func (service *QueryService) CreateAnnotations() Annotations {
	var res Annotations
	measureMap := make(map[string]Annotation)
	dimensionMap := make(map[string]Annotation)
	timeDimensionMap := make(map[string]Annotation)

	for _, measure := range service.Query.Measures {
		measureMap[measure] = Annotation{Title: GetTitle(measure), ShortTitle: GetShortTitle(measure), Type: GetMeasureType(measure)}
	}
	for _, dim := range service.Query.Dimensions {
		dimensionMap[dim] = Annotation{Title: GetTitle(dim), ShortTitle: GetShortTitle(dim), Type: GetDimensionType(dim)}
	}
	for _, timeD := range service.Query.TimeDimensions {
		timeDimensionMap[timeD.Dimension] = Annotation{Title: GetTitle(timeD.Dimension), ShortTitle: GetShortTitle(timeD.Dimension), Type: "time"}
	}
	res.Measures = measureMap
	res.Dimensions = dimensionMap
	res.TimeDimension = timeDimensionMap
	return res
}
