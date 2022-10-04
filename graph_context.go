package graph_engine_go

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
)

type GraphContext struct {
	graphConfig  *GraphConfig
	graphNodes   map[string]*GraphOperator
	inputNodes   []string
	computeNodes []string
	outputNodes  []string
	allNodes     []string
}

func NewGraphContext() *GraphContext {
	ctx := new(GraphContext)
	ctx.graphNodes = make(map[string]*GraphOperator, 16)
	return ctx
}

func (ctx *GraphContext) BuildGraph(graphConfig *GraphConfig) error {
	ctx.graphConfig = graphConfig
	// 创建所有节点
	nodeEmitMap := make(map[string]mapset.Set, 16)
	nodeDepMap := make(map[string]mapset.Set, 16)
	for _, nodeconfig := range graphConfig.Nodes {
		nodeName := nodeconfig.Name
		typeName := nodeconfig.NodeType
		// 创建对象
		node := CreateInstance(typeName)
		// 设置节点配置
		node.SetConfig(&nodeconfig)
		// 记录节点
		if _, ok := ctx.graphNodes[nodeName]; ok {
			return fmt.Errorf("not allow duplicated node name: %s", nodeName)
		}
		ctx.graphNodes[nodeName] = &node
		// 设置节点发布数据的信息
		nodeEmitMap[nodeName] = mapset.NewSet()
		for _, emitName := range nodeconfig.Emitters {
			if nodeEmitMap[nodeName].Contains(emitName) {
				return fmt.Errorf("not allow duplicated emitter name '%s' in node '%s'", emitName, nodeName)
			}
			nodeEmitMap[nodeName].Add(emitName)
		}
		// 设置节点依赖数据的信息
		nodeDepMap[nodeName] = mapset.NewSet()
		for _, depName := range nodeconfig.Depends {
			if nodeDepMap[nodeName].Contains(depName) {
				return fmt.Errorf("not allow duplicated depend name '%s' in node '%s'", depName, nodeName)
			}
			nodeDepMap[nodeName].Add(depName)
		}
		// 判断节点类型
		if len(nodeconfig.Emitters) == 0 {
			ctx.outputNodes = append(ctx.outputNodes, nodeName)
		} else if len(nodeconfig.Depends) == 0 {
			ctx.inputNodes = append(ctx.inputNodes, nodeName)
		} else {
			ctx.computeNodes = append(ctx.computeNodes, nodeName)
		}
	}
	fmt.Printf("num_inputs:%d, num_compute:%d, num_output:%d\r\n",
		len(ctx.inputNodes), len(ctx.computeNodes), len(ctx.outputNodes))
	// 找出所有对输出有贡献的节点
	var nodeFound []string
	validNodes := mapset.NewSet()
	for _, name := range ctx.outputNodes {
		nodeFound = append(nodeFound, name)
		validNodes.Add(name)
	}
	for idx := 0; idx < len(nodeFound); idx++ {
		name := nodeFound[idx]
		depends := nodeDepMap[name]
		for v := range depends.Iter() {
			// fmt.Printf("node %s, depend: %s\r\n", name, v)
			// 查找当前依赖数据的发布节点
			for emitName, c := range nodeEmitMap {
				if c.Contains(v) {
					if !validNodes.Contains(emitName) {
						validNodes.Add(emitName)
						nodeFound = append(nodeFound, emitName)
					}
					break
				}
			}
		}
	}
	ctx.allNodes = nodeFound
	fmt.Printf("found nodes: %v\r\n", ctx.allNodes)
	// 记录所有需要运行的节点
	ctx.inputNodes = ctx.inputNodes[:0]
	ctx.computeNodes = ctx.computeNodes[:0]
	ctx.outputNodes = ctx.outputNodes[:0]
	for _, name := range nodeFound {
		// 不依赖数据的节点为输入节点
		if nodeDepMap[name].Cardinality() == 0 {
			ctx.inputNodes = append(ctx.inputNodes, name)
		} else if nodeEmitMap[name].Cardinality() == 0 {
			ctx.outputNodes = append(ctx.outputNodes, name)
		} else {
			ctx.computeNodes = append(ctx.computeNodes, name)
		}
	}
	fmt.Printf("input:%v, compute:%v, output:%v\r\n", ctx.inputNodes, ctx.computeNodes, ctx.outputNodes)
	return nil
}
