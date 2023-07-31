package tests

import (
	"dashboard/internal/database"
	blockService "dashboard/internal/services/block"
	noSqlQuery "dashboard/internal/services/noSqlQuery"
	query "dashboard/internal/services/query"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

// Need to connect to a test database
func connectDb() {
	data := database.DatabaseInfo{DbHost: "0.0.0.0", DbPort: "27017", DbUsername: "root", DbPass: "pass12345", DbName: "test"}
	mongoDb := database.GetMongoDatabase()
	mongoDb.ConnectDatabase(data)
	database.SetMongoDatabase(mongoDb)
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

func getJoinObject() blockService.Join {
	join := blockService.Join{
		Name:         "Movies",
		LocalField:   "movie_id",
		ForeignField: "_id",
		Relationship: "one_to_one"}
	return join
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

func buildBlockQuery(name string) noSqlQuery.BlockQuery {
	blockQuery := noSqlQuery.BlockQuery{
		Measure:    []string{"count"},
		Dimensions: []string{"category", "time"},
		Name:       name,
	}
	return blockQuery
}

func buildAllBlockQueries() []noSqlQuery.BlockQuery {
	return []noSqlQuery.BlockQuery{buildBlockQuery("Stories"), buildBlockQuery("Movies")}
}

func TestQuery(t *testing.T) {
	connectDb()
	q := getQueryObject()

	t.Run("ParseQuery", func(t *testing.T) {
		res, err := noSqlQuery.ParseQuery(q)
		if err != nil {
			t.Fatalf("Err -> error during execution: %v", err)
		}
		log.Println(res)
	})
}

func TestBuildGroupStage(t *testing.T) {
	q := getQueryObject()
	j := getJoinObject()
	t.Run("Correct build stage", func(t *testing.T) {
		res := noSqlQuery.GenerateGroupStage(q.Dimensions, q.Measures, &j)
		log.Println(res)
		s := res["$group"].(bson.M)
		if s["_id"].(bson.M)["release_date"] != "$Movies.release_date" {
			t.Fatalf("Err -> \nWant %q\nGot %q", "$Movies.release_date", s["_id"].(bson.M)["release_date"])
		}
	})
}

func TestFindJoin(t *testing.T) {
	dimensions := []string{
		"Stories.category",
		"Movies.release_date",
	}
	blockName := noSqlQuery.FindBlockWithJoin(dimensions)
	log.Println(blockName)
	if blockName == nil {
		t.Fatalf("Err -> \nReturned nil")
	}
}
