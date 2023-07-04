package tests

import (
	"dashboard/internal/database"
	blockService "dashboard/internal/services/block"
	query "dashboard/internal/services/query"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
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

func buildStoriesBlock() blockService.BlockData {
	block := blockService.BlockData{
		Name: "Stories",
		Sql:  "Stories",
		Joins: []blockService.Join{
			{Name: "Movies", LocalField: "Stories.movie_id", ForeignField: "Movies.id", Relationship: "one_to_one"},
		},
		Measures: []blockService.Measures{
			{Name: "count", Sql: "_id", Type: "count"},
		},
		Dimensions: []blockService.Dimensions{
			{Name: "category", Sql: "category", Type: "string"},
			{Name: "isDraft", Sql: "isDraft", Type: "boolean"},
			{Name: "time", Sql: "time", Type: "string"},
			{Name: "movieId", Sql: "id", Type: "string", PrimaryKey: true},
		},
	}
	return block
}

func buildBlockQuery(name string) query.BlockQuery {
	blockQuery := query.BlockQuery{
		Measure:    []string{"count"},
		Dimensions: []string{"category", "time"},
		Name:       name,
	}
	return blockQuery
}

func buildAllBlockQueries() []query.BlockQuery {
	return []query.BlockQuery{buildBlockQuery("Stories"), buildBlockQuery("Movies")}
}

// func TestQuery(t *testing.T) {
// 	connectDb()
// 	q := getQueryObject()

// 	t.Run("ParseQuery", func(t *testing.T) {
// 		res, err := query.ParseQuery(q)
// 		if err != nil {
// 			t.Fatalf("Err -> error during execution: %v", err)
// 		}
// 		log.Println(res)
// 	})
// }

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

func TestBuildGroupStage(t *testing.T) {
	q := getQueryObject()

	t.Run("Correct build stage", func(t *testing.T) {
		res, err := query.BuildGroupStageFromDimensions(q.Dimensions)
		log.Println(res)
		if err != nil {
			t.Fatalf("Err -> \nReturned error: %v", err)
		}
		s := res["$group"].(bson.M)
		if s["_id"].(bson.M)["Movies"] != "$movie_id" {
			t.Fatalf("Err -> \nWant %q\nGot %q", "$movie_id", s["Movies"])
		}
	})
}

func TestFindJoin(t *testing.T) {
	dimensions := []string{
		"Stories.category",
		"Movies.release_date",
	}
	blockName := query.FindBlockWithJoin(dimensions)
	log.Println(blockName)
	if blockName == nil {
		t.Fatalf("Err -> \nReturned nil")
	}
}
