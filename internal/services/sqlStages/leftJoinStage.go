package sqlStages

import (
	"dashboard/internal/services/block"
	"fmt"
	"log"
	"strings"

	"golang.org/x/exp/slices"
)

func GenerateJoinClause(path []string, graph *block.JoinGraph) string {
	var result strings.Builder

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
	return result.String()
}

func GetLeftJoinPath(startTableBlock *block.BlockData, targetTableNames []string, graph *block.JoinGraph) []string {

	for _, targetTable := range targetTableNames {
		if startVertex, found := graph.Vertices[startTableBlock.Name]; found {
			path, relationshipFound := graph.FindJoinPath(startVertex, targetTable)
			if !relationshipFound || !slices.Contains(path, startTableBlock.Name) {
				startVertex = graph.Vertices[targetTable]
				log.Printf("relationship not found trying other way around: %v", startTableBlock.Name)
				path, relationshipFound = graph.FindJoinPath(startVertex, startTableBlock.Name)
			}
			if relationshipFound {
				return path
			}
		}
	}
	return nil
}
