package main

import (
	"dashboard/internal/services/block"
	"fmt"
	"sync"
)

var lock = &sync.Mutex{}

type single struct {
	Blocks []*block.FileData
}

var singleInstance *single

func GetInstance() *single {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating single instance now.")
			data, err := block.ReadAllBlocks("./example_schema")
			if err != nil {
				fmt.Println("Error trying to read all blocks: ", err)
			} else {
				singleInstance = &single{}
				singleInstance.Blocks = data
			}
		}
	}
	return singleInstance
}
