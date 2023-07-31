package database

type DatabaseInfo struct {
	DbHost     string `json:"db_host"`
	DbPort     string `json:"db_port"`
	DbUsername string `json:"db_username"`
	DbPass     string `json:"db_pass"`
	DbName     string `json:"db_name"`
}

type IDatabase interface {
	BuildDatabaseUri(DatabaseInfo) string
	ConnectDatabase(DatabaseInfo) error
}

var mongoDb = &mongoDatabase{db: nil}
var postgresDb = &postgresDatabase{db: nil}

func GetCurrentDatabase() IDatabase {
	if mongoDb.db == nil {
		return postgresDb
	}
	return mongoDb
}

func GetMongoDatabase() *mongoDatabase {
	return mongoDb
}

func GetPostgresDatabase() *postgresDatabase {
	return postgresDb
}

func SetMongoDatabase(db *mongoDatabase) {
	mongoDb = db
}

func SetPostgresDatabase(db *postgresDatabase) {
	postgresDb = db
}
