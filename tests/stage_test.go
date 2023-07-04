package tests

import (
	"dashboard/internal/services/query"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestBuildAFilter(t *testing.T) {
	q := getQueryObject()

	t.Run("correct filters", func(t *testing.T) {
		res, err := query.BuildAllFilters(q.Filters)
		log.Println(res)
		if err != nil {
			t.Fatalf("Err -> \nReturned error: %v", err)
		}
		match := res[0]["$match"].(bson.M)

		if match["isDraft"] != false {
			t.Fatalf("Err -> \nWant %v\nGot %q", false, match["isDraft"])
		}
	})
}

func TestTimeDimensions(t *testing.T) {
	q := getQueryObject()

	t.Run("correct time", func(t *testing.T) {
		res, err := query.BuildAllTimeDimensions(q.TimeDimensions)
		if err != nil {
			t.Fatalf("Err -> \nReturned error: %v", err)
		}
		match := res[0]["$match"].(bson.M)
		if match["time"].(bson.M)["$gte"] != "2015-01-01" {
			t.Fatalf("Err -> \nWant %q\nGot %q", "2015-01-01", match["$isDraft"])
		}

	})
}
