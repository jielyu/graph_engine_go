# graph_engine_go

***第一个golang项目，持续优化开发中...***

用于实现Golang模块的自动组装功能，便于模块解耦，易于更新维护。

**golang版本的图引擎**

在实际的业务生产环境中，经常需要对各种逻辑进行组合以实现业务需求。使用一套框架来完成各个业务逻辑的组装是比较便于维护的，既能节省研发人员的精力，又能保持运行稳定。

本仓库为Golang工程提供图引擎框架，也就是**以图的形式来组织程序中的数据和操作。** 一个操作总是依赖一些数据或者发射一些数据，包括输入和输出操作，因此这些任务可以使用图来描述，图中包括一些节点(操作)以及这些节点之间的连接关系(数据)。 图引擎是由数据驱动的。如果一个操作所依赖的数据都准备就绪，那么该操作就会自动运行起来。


## Install

在自己的工程根目录下使用如下指令

```shell
go get github.com/jielyu/graph_engine_go
```

在需要使用的文件中导入`graph_engine_go`

```go
import (
    ge "github.com/jielyu/graph_engine_go"
)
```

## User Guide

示例: [example](./example/README.md)

### 第1步 实现接口

实现[`ge.GraphOperator`](./graph_operator.go) 接口

```go
type GraphOperator interface {
	SetConfig(config *GraphNodeConfig) error
	SetUp(ctx *GraphContext) error
	Initailize(ctx *GraphContext) error
	Process(ctx *GraphContext) error
}
```

图引擎中的节点分为3种：

* 输入节点，不依赖其他节点提供的数据，每次运行时自动处于就绪状态。[输入节点示例](./example/nodes/two_number_generator.go)
* 中间计算节点，依赖其他节点提供数据，也会提供数据给其他节点依赖，需要等待依赖的数据全部就绪才能运行。[中间节点示例](./example/nodes/two_number_add_operator.go)
* 输出节点，依赖其他节点提供数据，但不为其他节点提供依赖数据，需要等待依赖的数据全部就绪才能运行。[输出节点示例](./example/nodes/number_printer.go)

**建议：每个节点单独放在一个文件，便于维护**

### 第2步 注册类型

在`init`函数中注册所定义的节点类型，例如

```go
func init() {
	ge.Register[TwoNumbersGeneratorOp]()
}
```

**建议：init函数与节点实现放在同一个文件中，便于解耦合**

### 第3步 编写配置

编写图结构的配置文件，可以参照例子 [main_graph](./example/config/main_graph.json)

图结构配置字段说明

```
"name": string, 指定图的名字
"type": string, 指定图的类型，[MainGraph/SubGraph]
"include": []string, 指定依赖的子图文件
"pool_size": int, 指定资源池的大小，越大的pool响应并发的能力越强，但消耗资源也越大
"num_threads": int, 指定单个context运行的线程数，暂时无用
"nodes": []GraphNodeConfig, 存储每个节点的配置信息
```

节点配置信息说明

```
{
    "name": string, 节点名字
    "type": string, 节点类型名，要求必须是注册过的类型名，否则报错
    "depends": map[string]string, 依赖的数据 
    {
        "Name": "data", Name是在程序中使用的，data是在当前配置文件内使用的，可以多个
    },
    "emitters": map[string]string, 发布的数据 
    {
        "Name": "data", Name是在程序中使用的，data是在当前配置文件内使用的，可以多个
    },
    "config": map[string]string, 用于设置一些配置项，需要自行在实现时解析
}
```

### 第4步 运行图引擎

准备工作就绪之后就可以把图引擎运行起来了，可以参照[示例](./example/main.go)

```go
package main

import (
	"fmt"
	"sync"

	_ "example/nodes"
	ge "github.com/jielyu/graph_engine_go"
)

var wg sync.WaitGroup

func main() {
	fmt.Println("welcome to graph_engine_go example")
	ge := ge.NewGraphEngine("./config/main_graph.json")
	times := 10
	wg.Add(times)
	for i := 0; i < times; i++ {
		go func(reqId uint64) {
			var v interface{}
			ge.Process(v, reqId)
			wg.Done()
		}(uint64(i))
	}
	wg.Wait()
}

```


## Coming Soon

* 子图机制，使用子图提高图配置的复用率
* 动态图机制，使用节点发射数据控制图的运行流程

