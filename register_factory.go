// 用于注册算子，并提供创建对象的接口
package graph_engine_go

import "fmt"

var (
	factory_registry = make(map[string]func() GraphOperator)
)

func RegisterClass(name string, fac_func func() GraphOperator) {
	factory_registry[name] = fac_func
}

func CreateInstance(name string) GraphOperator {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if f, ok := factory_registry[name]; ok {
		return f()
	} else {
		info_str := fmt.Sprintf("not found Class[%s]", name)
		panic(info_str)
	}
}
