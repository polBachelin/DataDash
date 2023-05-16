package block

import (
	"fmt"
	"sync"
)

var lock = &sync.Mutex{}

type single struct {
	Blocks []*FileData
}

var singleInstance *single

func GetInstance() *single {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating single instance now.")
			data, err := ReadAllBlocks("./schema")
			if err != nil {
				fmt.Println("Error trying to read all blocks: ", err)
			} else {
				singleInstance = &single{}
				singleInstance.Blocks = data
			}
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singleInstance
}
