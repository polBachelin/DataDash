package block

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var SCHEMA_PATH = "./schema/"

type FileData struct {
	Blocks []BlockData `yaml:"blocks"`
}

type BlockData struct {
	Name       string       `yaml:"name"`
	Sql        string       `yaml:"sql"`
	Measures   []Measures   `yaml:"measures"`
	Dimensions []Dimensions `yaml:"dimensions"`
}

type Measures struct {
	Name string `yaml:"name"`
	Sql  string `yaml:"sql"`
	Type string `yaml:"type"`
}

type Dimensions struct {
	Name string `yaml:"name"`
	Sql  string `yaml:"sql"`
	Type string `yaml:"type"`
}

func ReadFile(filename string) (*FileData, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	c := &FileData{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("Error in file %s: %v", filename, err)
	}
	return c, err
}
