// 用于注册算子，并提供创建对象的接口
package graph_engine_go

import "fmt"

var (
	factoryRegistry = make(map[string]func() GraphOperator)
)

func RegisterClass(name string, fac_func func() GraphOperator) {
	factoryRegistry[name] = fac_func
}

func CreateInstance(name string) (GraphOperator, error) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }()
	if f, ok := factoryRegistry[name]; ok {
		return f(), nil
	}
	return nil, fmt.Errorf("not found Class[%s]", name)
}
