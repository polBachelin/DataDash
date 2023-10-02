package sqlStages

import (
	"dashboard/internal/services/block"
	"fmt"
	"log"
	"strings"

	"golang.org/x/exp/slices"
)

func GenerateJoinClause(path []string, graph *block.JoinGraph) string {
	var joins strings.Builder

	for i := 0; i < len(path)-1; i++ {
		fromVertex := graph.Vertices[path[i]]
		toVertex := graph.Vertices[path[i+1]]

		joinParent, err := block.GetBlockJoinFromName(toVertex.Val.Name, fromVertex.Val)
		if err != nil {
			joinParent, _ = block.GetBlockJoinFromName(fromVertex.Val.Name, toVertex.Val)
		}
		joins.WriteString(fmt.Sprintf(" LEFT JOIN %s as %s ON %s.%s = %s.%s", toVertex.Val.Table, toVertex.Val.Name, toVertex.Val.Name, joinParent.LocalField, fromVertex.Val.Name, joinParent.ForeignField))
	}
	log.Printf("PATH FOR JOIN : %s", path)
	return joins.String()
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
