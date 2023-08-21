package query

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

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

// TODO: need to cut this function it is too long
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
	if len(selectStage) > 1 {
		sqlQuery.WriteString(service.Query.GenerateGroupByStage(len(selectStage)))
	}
	if len(service.Query.Filters) > 1 {
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
