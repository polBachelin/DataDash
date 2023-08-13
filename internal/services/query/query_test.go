package query

import (
	"dashboard/internal/database"
	"dashboard/internal/services/block"
	"log"
	"testing"
)

// Need to connect to a test database
func connectDb() bool {
	data := database.DatabaseInfo{DbHost: "172.24.0.1", DbPort: "5438", DbUsername: "postgres", DbPass: "postgres", DbName: "postgres"}
	postgres := database.GetPostgresDatabase()
	err := postgres.ConnectDatabase(data)
	if err != nil {
		log.Println("Error connecting database: ", err)
		return false
	}
	database.SetPostgresDatabase(postgres)
	return true
}

func getQueryObject() Query {
	q := Query{}
	q.Measures = []string{}
	q.Dimensions = []string{"Sale.amount"}
	f := Filter{Member: "Sale.amount", Operator: "gt", Values: []string{"9000"}}
	q.Filters = []Filter{f}
	timeDimension := TimeDimension{Dimension: "Sale.date", DateRange: []string{"2019-07-04", "2019-09-22"}, Granularity: "week"}
	q.TimeDimensions = []TimeDimension{timeDimension}
	q.Limit = 100
	q.Offset = 0
	q.Order = Order{}
	return q
}

func getQueryService() *QueryService {
	b := block.GetInstance().Blocks
	service := NewQueryService(getQueryObject(), database.GetCurrentDatabase(), block.NewGraph(b))
	return service
}

func TestSqlGeneration(t *testing.T) {
	service := getQueryService()
	if !connectDb() {
		t.Fatalf("Could not connect to db")
	}
	t.Run("SqlGeneration", func(t *testing.T) {
		json, _ := service.ParseQuery()
		log.Println(json)
	})
}
