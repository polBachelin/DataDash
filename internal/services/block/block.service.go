package block

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

var SCHEMA_PATH = "./schema/"

type FileData struct {
	Blocks []BlockData `yaml:"blocks"`
}

type BlockData struct {
	Name       string       `yaml:"name"`
	Table      string       `yaml:"table"`
	Joins      []Join       `yaml:"joins"`
	Measures   []Measures   `yaml:"measures"`
	Dimensions []Dimensions `yaml:"dimensions"`
}

type Join struct {
	Name         string `yaml:"name"`
	LocalField   string `yaml:"local_field"`
	ForeignField string `yaml:"foreign_field"`
	Relationship string `yaml:"relationship"`
}

type Measures struct {
	Name string `yaml:"name"`
	Sql  string `yaml:"sql"`
	Type string `yaml:"type"`
}

type Dimensions struct {
	Name       string `yaml:"name"`
	Sql        string `yaml:"sql"`
	Type       string `yaml:"type"`
	PrimaryKey bool   `yaml:"primary_key"`
}

func NewJoin() *Join {
	return &Join{}
}

func GetBlockFromName(name string) *BlockData {
	blockInstance := GetInstance()
	for _, blockData := range blockInstance.Blocks {
		b := slices.IndexFunc(blockData.Blocks, func(data BlockData) bool { return data.Name == name })
		if b != -1 {
			return &blockData.Blocks[b]
		}
	}
	return nil
}

func GetBlockJoinFromName(name string, block *BlockData) (Join, error) {
	join := slices.IndexFunc(block.Joins, func(data Join) bool { return data.Name == name })
	if join != -1 {
		return block.Joins[join], nil
	}
	return Join{}, errors.New("No join found for name: " + name)
}

func GetAllBlockNamesDifferent(compare string, members []string) []string {
	var res []string

	for _, member := range members {
		blockName := GetBlockName(member)
		if blockName != compare {
			b := GetBlockFromName(blockName)
			res = append(res, b.Name)
		}
	}
	return res
}

func GetBlockName(dimension string) string {
	parts := strings.Split(dimension, ".")
	return parts[0]
}

func GetBlockThatHasJoin(name string) *BlockData {
	blockInstance := GetInstance().Blocks

	for _, fileData := range blockInstance {
		for _, block := range fileData.Blocks {
			for _, blockJoin := range block.Joins {
				if blockJoin.Name == name {
					return &block
				}
			}
		}
	}
	return nil
}

func GetMeasureFromBlock(blockData *BlockData, measure string) (*Measures, error) {
	measureIndex := slices.IndexFunc(blockData.Measures, func(m Measures) bool { return m.Name == measure })
	if measureIndex == -1 {
		return nil, fmt.Errorf("no measure found for block %s", blockData.Name)
	}
	return &blockData.Measures[measureIndex], nil
}
