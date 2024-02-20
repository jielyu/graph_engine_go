package graph_engine_go

import (
	"fmt"
	"sync"
	"time"
)

type GraphEngine struct {
	graphConfig *GraphConfig    // 图结构配置
	contextPool []*GraphContext // 运行池
}

// 创建一个图引擎
func NewGraphEngine(jsonFile string) *GraphEngine {
	ge := new(GraphEngine)
	graphConfig, err := LoadGraphConfig(jsonFile)
	if err != nil {
		panic(err)
	}
	ge.graphConfig = graphConfig
	if ge.graphConfig.PoolSize < 0 {
		ge.graphConfig.PoolSize = 1
	}
	// 创建上下文运行池
	ge.contextPool = make([]*GraphContext, 0, ge.graphConfig.PoolSize)
	for i := 0; i < ge.graphConfig.PoolSize; i++ {
		ctx := NewGraphContext()
		ctx.Id = i
		ctx.Busy = false
		err := ctx.Build(ge.graphConfig)
		if err != nil {
			panic(err)
		}
		ge.contextPool = append(ge.contextPool, ctx)
	}
	return ge
}

var poolMutex sync.Mutex

func (ge *GraphEngine) selectIdleCtx() *GraphContext {
	poolMutex.Lock()
	defer poolMutex.Unlock()
	for _, v := range ge.contextPool {
		if !v.Busy {
			v.Busy = true
			return v
		}
	}
	return nil
}

// 用于执行一次处理任务
// inputData可以在输入节点强制转换为对应结构
// reqId用于标识每次处理，最好保证互异
func (ge *GraphEngine) Process(inputData interface{}, reqId uint64) (interface{}, error) {
	// 选择空闲的context，否则等待
	ctx := ge.selectIdleCtx()
	for ctx == nil {
		time.Sleep(time.Duration(1) * time.Microsecond)
		ctx = ge.selectIdleCtx()
	}
	// 清楚busy标志
	defer func() {
		ctx.Busy = false
	}()
	fmt.Printf("[reqId:%d] select pool %d to process\r\n", reqId, ctx.Id)

	ctx.InputData = inputData
	ctx.ReqId = reqId
	err := ctx.Process()
	// fmt.Printf("ctx %d run finished.addr:%p\r\n", ctx.Id, ctx)
	return ctx.OutputData, err
}
