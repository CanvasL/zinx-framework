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


// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	fmt.Println("receive from client, msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))
	if err := request.GetConnection().SendMsg(1, []byte("This is China!This is China!This is China!")); err != nil {
		fmt.Println("send msg failed, ", err)
	}
}


func main() {
	// 1.创建一个句柄(使用zinx的api)
	s := znet.NewServer("[zinx V0.5]")
	// 2.给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	// 3.启动server
	s.Serve()
}