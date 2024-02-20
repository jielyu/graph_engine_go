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
	times := 1000
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
