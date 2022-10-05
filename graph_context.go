package graph_engine_go

import (
	"fmt"
	"sync/atomic"
	"time"

	mapset "github.com/deckarep/golang-set"
)

type GraphContext struct {
	graphConfig    *GraphConfig
	graphNodes     map[string]GraphOperator // 存储所有节点对象
	nodeEmitMap    map[string]mapset.Set    // 存储每个节点发布的数据名字
	nodeDepMap     map[string]mapset.Set    // 存储每个节点依赖数据的名字
	inputNodes     []string                 // 存储所有需要运行的输入节点的名字
	computeNodes   []string                 // 存储所有需要运行的计算节点的名字
	outputNodes    []string                 // 存储所有需要运行的输出节点名字
	allNodes       []string                 // 存储所有需要运行的节点名字
	allGraphData   map[string]*GraphData    // 存储所有发布的数据
	allGraphDep    map[string]*GraphDep     // 存储所有依赖的数据结构
	dataByNodeMap  map[string]string        // 用于存储发布数据与发布节点之间的映射
	dataForNodeMap map[string][]string      // 用于存储发布的数据与依赖该数据的节点之间的映射
}

// 创建GraphContext对象
func NewGraphContext() *GraphContext {
	ctx := new(GraphContext)
	ctx.graphNodes = make(map[string]GraphOperator, 16)
	ctx.allGraphData = make(map[string]*GraphData, 32)
	ctx.allGraphDep = make(map[string]*GraphDep, 32)
	ctx.dataByNodeMap = make(map[string]string, 32)
	ctx.dataForNodeMap = make(map[string][]string, 32)
	return ctx
}

// 用于创建发布数据的对象
func (ctx *GraphContext) NameEmit(name string) *GraphData {
	// 如果发布的数据已经存在，则直接返回。只可能是NameDepend()中发布的，不会重复
	if _, ok := ctx.allGraphData[name]; ok {
		return ctx.allGraphData[name]
	}
	// 创建发布数据结构
	graphData := new(GraphData)
	graphData.Name = name
	graphData.Data = nil
	ctx.allGraphData[name] = graphData
	fmt.Printf("create GraphData '%s' succ\r\n", name)
	return graphData
}

// 用于创建依赖数据的对象
func (ctx *GraphContext) NameDepend(name string) *GraphDep {
	// 依赖数据已经存在的情况下，直接返回
	if _, ok := ctx.allGraphDep[name]; ok {
		return ctx.allGraphDep[name]
	}
	// 依赖尚未发布的数据，先操作发布，创建时已经避免重复，这里不需检查
	if _, ok := ctx.allGraphData[name]; !ok {
		ctx.NameEmit(name)
	}
	// 创建依赖数据结构
	graphDep := new(GraphDep)
	graphDep.RefGraphData = ctx.allGraphData[name]
	ctx.allGraphDep[name] = graphDep
	return graphDep
}

// 创建图结构
func (ctx *GraphContext) Build(graphConfig *GraphConfig) error {
	// 创建节点
	err := ctx.create(graphConfig)
	if err != nil {
		return err
	}
	// 创建数据
	err = ctx.setup()
	if err != nil {
		return err
	}
	// 初始化
	err = ctx.initailize()
	if err != nil {
		return err
	}
	return nil
}

// 运行图结构
func (ctx *GraphContext) Process() error {
	// 记录节点依赖情况
	nodeState := make(map[string]*int32)
	for _, name := range ctx.allNodes {
		var stat int32 = int32(ctx.nodeDepMap[name].Cardinality())
		nodeState[name] = &stat
	}
	readyNode := make([]string, 0, len(nodeState))
	for len(nodeState) > 0 {
		// 查找准备就绪的节点
		for k, v := range nodeState {
			if *v == 0 {
				delete(nodeState, k)
				readyNode = append(readyNode, k)
			}
		}
		// 运行准备就绪的节点
		for _, name := range readyNode {
			go func(nname string) {
				fmt.Println("process node ", nname)
				op := ctx.graphNodes[nname]
				op.Process(ctx)
				for v := range ctx.nodeEmitMap[nname].Iter() {
					emitName := v.(string)
					for _, nodeName := range ctx.dataForNodeMap[emitName] {
						atomic.AddInt32(nodeState[nodeName], -1)
					}
				}
			}(name)
		}
		readyNode = readyNode[:0]
		time.Sleep(time.Duration(1) * time.Microsecond)
	}
	return nil
}

// 构建各个节点
func (ctx *GraphContext) create(graphConfig *GraphConfig) error {
	ctx.graphConfig = graphConfig
	// 创建所有节点
	nodeEmitMap := make(map[string]mapset.Set, 16)
	nodeDepMap := make(map[string]mapset.Set, 16)
	for _, nodeconfig := range graphConfig.Nodes {
		nodeName := nodeconfig.Name
		typeName := nodeconfig.NodeType
		// 创建对象
		node, err := CreateInstance(typeName)
		if err != nil {
			return err
		}
		// 设置节点配置
		fmt.Println(nodeName, nodeconfig)
		var configCopy GraphNodeConfig = nodeconfig
		node.SetConfig(&configCopy)
		// 记录节点
		if _, ok := ctx.graphNodes[nodeName]; ok {
			return fmt.Errorf("not allow duplicated node name: %s", nodeName)
		}
		ctx.graphNodes[nodeName] = node
		// 设置节点发布数据的信息
		nodeEmitMap[nodeName] = mapset.NewSet()
		for _, emitName := range nodeconfig.Emitters {
			// 避免同一个节点发布两份一样的数据
			if nodeEmitMap[nodeName].Contains(emitName) {
				return fmt.Errorf("not allow duplicated emitter name '%s' in node '%s'", emitName, nodeName)
			}
			nodeEmitMap[nodeName].Add(emitName)
			// 记录发布数据的节点，同时避免不同节点重复发布数据
			if _, ok := ctx.dataByNodeMap[emitName]; ok {
				return fmt.Errorf("not allow duplicated emitter name '%s' in node '%s' and '%s'",
					emitName, nodeName, ctx.dataByNodeMap[emitName])
			}
			ctx.dataByNodeMap[emitName] = nodeName
		}
		// 设置节点依赖数据的信息
		nodeDepMap[nodeName] = mapset.NewSet()
		for _, depName := range nodeconfig.Depends {
			if nodeDepMap[nodeName].Contains(depName) {
				return fmt.Errorf("not allow duplicated depend name '%s' in node '%s'", depName, nodeName)
			}
			nodeDepMap[nodeName].Add(depName)
			// 记录依赖数据的所有节点名字
			ctx.dataForNodeMap[depName] = append(ctx.dataForNodeMap[depName], nodeName)
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
	ctx.nodeEmitMap = nodeEmitMap
	ctx.nodeDepMap = nodeDepMap
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
	fmt.Printf("input:%v, compute:%v, output:%v\r\n",
		ctx.inputNodes, ctx.computeNodes, ctx.outputNodes)
	return nil
}

// 调用每一个node的SetUp函数
func (ctx *GraphContext) setup() error {
	for _, name := range ctx.allNodes {
		fmt.Printf("start to setup node '%s'\r\n", name)
		err := ctx.graphNodes[name].SetUp(ctx)
		if err != nil {
			fmt.Printf("failed to setup node '%s'", name)
			return err
		}
	}
	return nil
}

// 调用每一个node的Initialize函数
func (ctx *GraphContext) initailize() error {
	for _, name := range ctx.allNodes {
		err := ctx.graphNodes[name].Initailize(ctx)
		if err != nil {
			fmt.Printf("failed to initialize node '%s'", name)
			return err
		}
	}
	return nil
}
