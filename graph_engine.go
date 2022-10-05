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
	poolFlag    []int32         // 标志位
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
	ge.poolFlag = make([]int32, 0, ge.graphConfig.PoolSize)
	for i := 0; i < ge.graphConfig.PoolSize; i++ {
		ctx := NewGraphContext()
		ctx.Build(ge.graphConfig)
		ge.contextPool = append(ge.contextPool, ctx)
		var busy int32 = 0
		ge.poolFlag = append(ge.poolFlag, busy)
	}
	return ge
}

var selectLocker sync.Mutex

func (ge *GraphEngine) selectIdleCtx() int {
	selectLocker.Lock()
	defer selectLocker.Unlock()
	for idx, v := range ge.poolFlag {
		if v == 0 {
			ge.poolFlag[idx] = 1
			return idx
		}
	}
	return -1
}

func (ge *GraphEngine) Process() error {
	ctxIdx := ge.selectIdleCtx()
	for ctxIdx == -1 {
		time.Sleep(time.Duration(1) * time.Microsecond)
		ctxIdx = ge.selectIdleCtx()
		// fmt.Printf("ctxIdx:%d", ctxIdx)
	}
	fmt.Println("select pool ", ctxIdx, " to process")
	ge.contextPool[ctxIdx].Process()
	ge.poolFlag[ctxIdx] = 0
	return nil
}
