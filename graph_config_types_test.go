package graph_engine_go

import "testing"
import "fmt"
import "encoding/json"

func TestGraphConfig(t *testing.T) {
	fmt.Println("Testing for Graph Config")
	// 从文件中读取json
	var jsonPath string = "./config/test/graph_test.json"
	graphConfig, err := LoadGraphConfig(jsonPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load graph config, err:%v", err))
	}
	// 序列化为json数据
	jsonStr, err := json.Marshal(graphConfig)
	if err != nil {
		panic("Failed convert GraphConfig to json")
	}
	fmt.Printf("%s\n", jsonStr)
	fmt.Println(graphConfig.Nodes[0].Emitters[0].Name == "")
}
