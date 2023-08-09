package query

import (
	"dashboard/internal/database"
	"dashboard/internal/services/block"
	"testing"
)

// Need to connect to a test database
func connectDb() {
	data := database.DatabaseInfo{DbHost: "172.24.0.1", DbPort: "5438", DbUsername: "postgres", DbPass: "postgres", DbName: "ecom"}
	mongoDb := database.GetMongoDatabase()
	mongoDb.ConnectDatabase(data)
	database.SetMongoDatabase(mongoDb)
}

func getQueryObject() Query {
	q := Query{}
	q.Measures = []string{"Sale.count"}
	q.Dimensions = []string{"Status_name.status_name"}
	f := Filter{}
	q.Filters = []Filter{f}
	timeDimension := TimeDimension{}
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

func TestSqlGeneratoin(t *testing.T) {
	service := getQueryService()
	t.Run("SqlGeneration", func(t *testing.T) {
		service.ParseQuery()
	})
}
