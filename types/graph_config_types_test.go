package types_test

import "testing"
import "fmt"
import "encoding/json"
import "github.com/jielyu/graph_engine_go/types"

func TestGraphConfig(t *testing.T) {
	fmt.Println("Testing for Graph Config")
	// 从文件中读取json
	var jsonPath string = "../config/graph_test.json"
	graphConfig, err := types.LoadGraphConfig(jsonPath)
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
