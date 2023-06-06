package tests

import (
	"dashboard/internal/database"
	query "dashboard/internal/services/query"
	"log"
	"testing"

	"golang.org/x/exp/slices"
)

// Need to connect to a test database
func connectDb() {
	data := database.DatabaseInfo{DbHost: "0.0.0.0", DbPort: "27017", DbUsername: "root", DbPass: "pass12345", DbName: "test"}
	database.ConnectDatabase(data)
}

func getQueryObject() query.Query {
	q := query.Query{}
	q.Measures = []string{"Stories.count"}
	q.Dimensions = []string{"Stories.category", "Stories.time", "Movies.release_date"}
	f := query.Filter{Member: "Stories.isDraft", Operator: "equals", Values: []string{"No"}}
	q.Filters = []query.Filter{f}
	timeDimension := query.TimeDimension{Dimension: "Stories.time", DateRange: []string{"2015-01-01", "2015-12-31"}, Granularity: "day"}
	q.TimeDimensions = []query.TimeDimension{timeDimension}
	q.Limit = 100
	q.Offset = 0
	q.Order = query.Order{DimensionName: []string{"Stories.time"}, DimensionOrder: []string{"asc"}, MeasureName: []string{"Stories.count"}, MeasureOrder: []string{"desc"}}
	return q
}

func TestQuery(t *testing.T) {
	connectDb()
	q := getQueryObject()

	t.Run("ParseQuery", func(t *testing.T) {
		res, err := query.ParseQuery(q)
		if err != nil {
			t.Fatalf("Err -> error during execution: %v", err)
		}
		log.Println(res)
		if res.Data[0].DimensionType != "category" {
			t.Errorf("Err -> \nWant %q\nGot %q", "category", res.Data[0].DimensionType)
		}
		if res.Data[0].Dimension != "Fiction" {
			t.Errorf("Err -> \nWant %q\nGot %q", "Fiction", res.Data[0].Dimension)
		}
		if res.Data[0].MeasureType != "count" {
			t.Errorf("Err -> \nWant %q\nGot %q", "count", res.Data[0].MeasureType)
		}
		if res.Data[0].Measure != "4" {
			t.Errorf("Err -> \nWant %q\nGot %q", "4", res.Data[0].Measure)
		}

	})
}

func TestBlockQuery(t *testing.T) {
	q := getQueryObject()

	t.Run("GetBlockQueriesFromQuery", func(t *testing.T) {
		res := query.GetBlockQueriesFromQuery(q)
		log.Println(res)
		if res[0].Name != "Stories" {
			t.Fatalf("Err -> \nWant %q\nGot %q", "Stories", res[0].Name)
		}
		if !slices.Contains(res[0].Dimensions, "category") {
			t.Fatalf("Err -> \nWant %q\nGot %q", "category", res[0].Dimensions)
		}
		if !slices.Contains(res[1].Dimensions, "release_date") {
			t.Fatalf("Err -> \nWant %q\nGot %q", "release_date", res[1].Dimensions)
		}
	})
}
