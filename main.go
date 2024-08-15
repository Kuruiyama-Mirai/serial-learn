package main

import (
	"Demo/App/routers"
	"fmt"
)

func main() {

	r := routers.InitRouter()
	//加载模版和静态资源

	//默认监听在0.0.0.0:8080
	err := r.Run()
	if err != nil {
		fmt.Printf("启动失败, %v", err)
		return
	}
}
