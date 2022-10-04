package graph_engine_go

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// 节点配置
type GraphNodeConfig struct {
	Name     string            `json:"name"`
	NodeType string            `json:"type"`
	Emitters map[string]string `json:"emitters"`
	Depends  map[string]string `json:"depends"`
	Config   map[string]string `json:"config"`
}

// 整图结构的配置
type GraphConfig struct {
	Name       string            `json:"name"`
	TypeName   string            `json:"type"`
	Include    []string          `json:"include"`
	PoolSize   int               `json:"pool_size"`
	NumThreads int               `json:"num_threads"`
	Nodes      []GraphNodeConfig `json:"nodes"`
}

// 载入关于图结构信息描述的json文件
func LoadGraphConfig(configJson string) (*GraphConfig, error) {
	// 检查路径是否存在
	_, err := os.Stat(configJson)
	if os.IsNotExist(err) {
		return nil, nil
	}
	// 打开并读取json文件
	jsonFile, err := os.Open(configJson)
	if err != nil {
		return nil, fmt.Errorf("Failed_to_parse_json_file:%s", configJson)
	}
	defer jsonFile.Close()
	byteValues, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("Failed_to_read_data_from_json_file:%s", configJson)
	}
	// 解析json数据
	var graphConfig GraphConfig
	err = json.Unmarshal(byteValues, &graphConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed_to_parse_json:%s", configJson)
	}
	return &graphConfig, nil
}
