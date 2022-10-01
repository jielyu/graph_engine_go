package graph_engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type EmitterConfig struct {
	Name string `json:"name,omitempty"`
	Data string `json:"data"`
}

type DependConfig struct {
	Name string `json:"name,omitempty"`
	Data string `json:"data"`
}

type GraphNodeConfig struct {
	Name     string          `json:"name"`
	Depends  []DependConfig  `json:"depends"`
	Emitters []EmitterConfig `json:"emitters"`
	TypeName string          `json:"typeName"`
}

type GraphConfig struct {
	Name  string            `json:"name"`
	Nodes []GraphNodeConfig `json:"nodes"`
}

/***
* Function: 用于从配置文件加载图结构配置
* Args:
*   configJson string, 配置文件路径
* Returns:
*   GraphConfig, 图结构配置
*   error,       错误信息
 */
func LoadGraphConfig(configJson string) (GraphConfig, error) {
	var graphConfig GraphConfig
	// 检查路径是否存在
	_, err := os.Stat(configJson)
	if os.IsNotExist(err) {
		return graphConfig, nil
	}
	// 打开并读取json文件
	jsonFile, err := os.Open(configJson)
	if err != nil {
		return graphConfig, fmt.Errorf("Failed_to_parse_json_file:%s", configJson)
	}
	defer jsonFile.Close()
	byteValues, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return graphConfig, fmt.Errorf("Failed_to_read_data_from_json_file:%s", configJson)
	}
	// 解析json数据
	err = json.Unmarshal(byteValues, &graphConfig)
	if err != nil {
		return graphConfig, fmt.Errorf("Failed_to_parse_json:%s", byteValues)
	}
	return graphConfig, nil
}
