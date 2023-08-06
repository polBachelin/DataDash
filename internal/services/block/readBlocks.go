package block

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//TODO: need to validate yaml file for required and optional fields

func ReadBlockFile(filename string) (*FileData, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	c := &FileData{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("error in file %s: %v", filename, err)
	}
	return c, err
}

func ReadAllBlocks(directory string) ([]*FileData, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("Error in directory: %v", err)
		return nil, err
	}
	data := make([]*FileData, 0, len(entries))
	for _, e := range entries {
		block, err := ReadBlockFile(filepath.Join(directory, e.Name()))
		if err != nil {
			return data, err
		}
		data = append(data, block)
	}
	return data, err
}
