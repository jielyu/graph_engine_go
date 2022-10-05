package graph_engine_go

import (
	"fmt"
	"sync"
	// "sync/atomic"
	"time"
)

type GraphEngine struct {
	graphConfig *GraphConfig
	contextPool []*GraphContext // 运行池
}

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
	ge.contextPool = make([]*GraphContext, 0, ge.graphConfig.PoolSize)
	for i := 0; i < ge.graphConfig.PoolSize; i++ {
		ctx := NewGraphContext()
		ctx.Id = i
		ctx.Busy = false
		ctx.Build(ge.graphConfig)
		ge.contextPool = append(ge.contextPool, ctx)
	}
	return ge
}

var selectLocker sync.Mutex

func (ge *GraphEngine) selectIdleCtx() *GraphContext {
	selectLocker.Lock()
	defer selectLocker.Unlock()
	for _, v := range ge.contextPool {
		if !v.Busy {
			v.Busy = true
			return v
		}
	}
	return nil
}

func (ge *GraphEngine) Process(inputData interface{}) (interface{}, error) {
	ctx := ge.selectIdleCtx()
	for ctx == nil {
		time.Sleep(time.Duration(1) * time.Microsecond)
		ctx = ge.selectIdleCtx()
	}
	// 清楚busy标志
	defer func() {
		ctx.Busy = false
	}()
	fmt.Println("select pool ", ctx.Id, " to process")
	ctx.InputData = inputData
	err := ctx.Process()
	return ctx.OutputData, err
}
