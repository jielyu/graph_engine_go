# graph_engine_go

***第一个golang项目，持续优化开发中...***

用于实现Golang模块的自动组装功能，便于模块解耦，易于更新维护。

## Introduction

**golang版本的图引擎**

在实际的业务生产环境中，经常需要对各种逻辑进行组合以实现业务需求。使用一套框架来完成各个业务逻辑的组装是比较便于维护的，既能节省研发人员的精力，又能保持运行稳定。

本仓库为Golang工程提供图引擎框架，也就是**以图的形式来组织程序中的数据和操作。** 一个操作总是依赖一些数据或者发射一些数据，包括输入和输出操作，因此这些任务可以使用图来描述，图中包括一些节点(操作)以及这些节点之间的连接关系(数据)。 图引擎是由数据驱动的。如果一个操作所依赖的数据都准备就绪，那么该操作就会自动运行起来。

## Features

* **以图的方式组织业务逻辑，满足复杂业务需求**

>组织业务逻辑的方式大体可以分两类，一类是完全使用代码组织，通常用于学习和试验性的工程；另一类是使用配置文件+固定结构的框架，通常用于生产环境下。对于第二类，**一种比较简单的实现就是线性结构**，按照配置文件中的顺序执行业务逻辑，缺点是不够灵活；**另一种就是本项目采取的图结构**，配置文件中的信息描述节点和依赖，框架按照图结构组装节点，可以方便地实现复杂业务逻辑。

* **业务逻辑解耦，便于维护升级**

>把业务逻辑拆分到每个节点中，可以逻辑模块进行更新和日常维护

* **自动并行执行没有依赖关系的所有逻辑**

>图结构的组装框架自动完成没有依赖关系的节点并行执行，最大程度提升程序运行效率

* **尽量保持内存复用**

>同一处理池在处理不同请求的时候，尽量复用最开始创建的内存。能够最大程度避免程序运行过程中出现内存剧烈波动。

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
	// 设置并记录节点的config信息
	SetConfig(config *GraphNodeConfig) error
	// 建立图结构，主要用于创建依赖和发布的数据结构，只在创建时执行一次
	SetUp(ctx *GraphContext) error
	// 初始化操作，一般用于设置初始化参数，只在创建时执行一次
	Initailize(ctx *GraphContext) error
	// 执行任务，业务逻辑主要封装在这个函数，每次都会执行
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


## Future

* 子图机制，使用子图提高图配置的复用率
* 动态图机制，使用节点发射数据控制图的运行流程

