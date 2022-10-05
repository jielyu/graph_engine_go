package graph_engine_go

import (
	"sync"
	"testing"
)

var wg sync.WaitGroup

func TestGraphEngine(t *testing.T) {
	ge := NewGraphEngine("./config/test/main_graph.json")
	times := 100
	wg.Add(times)
	for i := 0; i < times; i++ {
		go func() {
			ge.Process()
			wg.Done()
		}()
	}
	wg.Wait()
}
