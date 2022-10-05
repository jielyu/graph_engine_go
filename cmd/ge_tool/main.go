package main

import (
	"fmt"
	"github.com/jielyu/graph_engine_go"
	"reflect"
	// "github.com/urfave/cli/v2"
	// log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Welcome to Graph Engine")
	// ctx := graph_engine_go.NewGraphContext()
	typeOfCtx := reflect.TypeOf(graph_engine_go.GraphContext{})
	fmt.Println(typeOfCtx.Name())
}
