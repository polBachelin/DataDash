package tests

import (
	"dashboard/internal/database"
	query "dashboard/internal/services/query"
	"testing"
)

// Need to connect to a test database
func connectDb() {
	data := database.DatabaseInfo{DbHost: "0.0.0.0", DbPort: "27017", DbUsername: "root", DbPass: "pass12345", DbName: "test"}
	database.ConnectDatabase(data)
}

func TestQuery(t *testing.T) {
	connectDb()
	q := query.Query{}
	q.Measures = []string{"Stories.count"}
	q.Dimensions = []string{"Stories.category"}
	f := query.Filter{Member: "Stories.isDraft", Operator: "equals", Values: []string{"No"}}
	q.Filters = []query.Filter{f}
	timeDimension := query.TimeDimension{Dimension: "Stories.time", DateRange: []string{"2015-01-01", "2015-12-31"}, Granularity: "day"}
	q.TimeDimensions = []query.TimeDimension{timeDimension}
	q.Limit = 100
	q.Offset = 0
	q.Order = query.Order{DimensionName: []string{"Stories.time"}, DimensionOrder: []string{"asc"}, MeasureName: []string{"Stories.count"}, MeasureOrder: []string{"desc"}}

	t.Run("ParseQuery", func(t *testing.T) {
		res := query.ParseQuery(q)
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
