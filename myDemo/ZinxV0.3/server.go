package main

import (
	"fmt"
	"zinx_framework_demo/zinx/ziface"
	"zinx_framework_demo/zinx/znet"
)

/*
	基于zinx来开发的服务器端应用程序
 */

// ping test
type PingRouter struct {
	znet.BaseRouter
}

// Test PreHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("callback before ping failed, err:", err)
	}
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("callback ping...ping...ping failed, err:", err)
	}
}

// Test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("callback after ping failed, err:", err)
	}
}


func main() {
	// 1.创建一个句柄(使用zinx的api)
	s := znet.NewServer("[zinx V0.3]")
	// 2.给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	// 3.启动server
	s.Serve()
}