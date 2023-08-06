package block

import (
	"dashboard/pkg/utils"
	"log"
	"sync"
)

var lock = &sync.Mutex{}

type single struct {
	Blocks []*FileData
}

var singleInstance *single = nil

func GetInstance() *single {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			path := utils.GetEnvVar("SCHEMA_PATH", "./example_schema/sale_db")
			log.Println("Creating single instance now.")
			data, err := ReadAllBlocks(path)
			if err != nil {
				log.Println("Error trying to read all blocks: ", err)
			} else {
				singleInstance = &single{}
				singleInstance.Blocks = data
			}
		} else {
			log.Println("Single instance already created.")
		}
	} else {
		log.Println("Single instance already created.")
	}
	return singleInstance
}
