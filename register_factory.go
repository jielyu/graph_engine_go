// 用于注册算子，并提供创建对象的接口
package graph_engine_go

import (
	"fmt"
	"reflect"
)

var (
	factoryRegistry = make(map[string]func() GraphOperator)
)

// 注册一个类型
func RegisterClass(name string, fac_func func() GraphOperator) error {
	if _, ok := factoryRegistry[name]; ok {
		return fmt.Errorf("not allow duplicated struct name:%s", name)
	}
	fmt.Printf("register type: %s\r\n", name)
	factoryRegistry[name] = fac_func
	return nil
}

// 使用范型和反射机制简化注册函数
func Register[T any]() {
	var tmp T
	typeOfT := reflect.TypeOf(tmp)
	err := RegisterClass(typeOfT.Name(), func() GraphOperator {
		var node T
		return any(&node).(GraphOperator)
	})
	if err != nil {
		panic(err)
	}
}

// 根据类型名字创建对象
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
