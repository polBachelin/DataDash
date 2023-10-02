package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type postgresDatabase struct {
	db *sql.DB
}

func (x *postgresDatabase) BuildDatabaseUri(dbData DatabaseInfo) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbData.DbHost, dbData.DbPort, dbData.DbUsername, dbData.DbPass, dbData.DbName)
}

func (x *postgresDatabase) ConnectDatabase(dbData DatabaseInfo) error {
	uri := x.BuildDatabaseUri(dbData)
	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	err = db.Ping()
	log.Println(err)
	if err != nil {
		return err
	}
	x.db = db
	return nil
}

func (x *postgresDatabase) GetDatabaseConnection() *sql.DB {
	if x.db != nil {
		return x.db
	}
	fmt.Println("No database connection")
	return nil
}

func (x *postgresDatabase) ExecuteQuery(query interface{}) (interface{}, error) {
	log.Printf("%s", fmt.Sprintf("%s", query))
	res, err := x.db.Query(fmt.Sprintf("%s", query))
	if err != nil {
		fmt.Println("Error executing query: ", err)
		return nil, err
	}
	return res, nil
}

func (x *postgresDatabase) QueryResultToJson(rows interface{}) ([]map[string]interface{}, error) {
	sqlRows := rows.(*sql.Rows)
	columns, err := sqlRows.Columns()
	if err != nil {
		log.Println("Error in retrieving columns: ", err)
		return nil, err
	}
	values := make([]interface{}, len(columns))
	for i := range values {
		var v interface{}
		values[i] = &v
	}
	var resJson []map[string]interface{}
	for sqlRows.Next() {
		err := sqlRows.Scan(values...)
		if err != nil {
			log.Println("Error: ", err)
		}
		rowData := make(map[string]interface{})
		for i, colName := range columns {
			val := *values[i].(*interface{})
			switch v := val.(type) {
			case nil:
				rowData[colName] = 0
			case int64, int32, int16, int8, float64, float32:
				rowData[colName] = v
			case []byte:
				rowData[colName] = string(v)
			case time.Time:
				formatted := v.Format("2006-01-02T15:04:05.999")
				rowData[colName] = formatted
			default:
				rowData[colName] = val
			}
		}
		resJson = append(resJson, rowData)
	}
	if err = sqlRows.Err(); err != nil {
		return nil, err
	}
	defer sqlRows.Close()
	return resJson, nil
}
