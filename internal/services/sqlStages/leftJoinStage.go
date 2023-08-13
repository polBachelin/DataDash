package sqlStages

import (
	"dashboard/internal/services/block"
	"fmt"
	"strings"
)

func GenerateJoinClause(path []string, graph *block.JoinGraph) string {
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

func BuildLeftJoinSql(startTableBlock *block.BlockData, targetTableNames []string, graph *block.JoinGraph) string {
	var joins strings.Builder

	for _, targetTable := range targetTableNames {
		if startVertex, found := graph.Vertices[startTableBlock.Name]; found {
			path, relationshipFound := graph.FindJoinPath(startVertex, targetTable)
			if relationshipFound {
				joins.WriteString(GenerateJoinClause(path, graph))
			}
		}
	}
	return joins.String()
}
