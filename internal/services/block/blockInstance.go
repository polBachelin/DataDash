package block

import (
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
			log.Println("Creating single instance now.")
			data, err := ReadAllBlocks("./schema")
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
