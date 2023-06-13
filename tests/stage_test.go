package tests

import (
	"dashboard/internal/services/query"
	"log"
	"testing"
)

func TestBuildAFilter(t *testing.T) {
	q := getQueryObject()

	t.Run("correct filters", func(t *testing.T) {
		res, err := query.BuildAllFilters(q.Filters)
		log.Println(res)
		if err != nil {
			t.Fatalf("Err -> \nReturned error: %v", err)
		}
	})
}
