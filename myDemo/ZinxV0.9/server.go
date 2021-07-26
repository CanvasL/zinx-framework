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
	if err := request.GetConnection().SendMsg(200, []byte("ping, ping, ping!")); err != nil {
		fmt.Println("send msg failed, ", err)
	}
}

// hello test
type HelloRouter struct {
	znet.BaseRouter
}
// Test Handle
func (this *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	fmt.Println("receive from client, msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))
	if err := request.GetConnection().SendMsg(201, []byte("Hello, Hello, Hello!")); err != nil {
		fmt.Println("send msg failed, ", err)
	}
}

// 创建链接之后执行的钩子函数
func DoConnBegin(conn ziface.IConnection) {
	fmt.Println("===>DoConnBegin is called...")
	if err := conn.SendMsg(202, []byte(">>Hook: 创建链接后的hook函数被执行了！")); err != nil {
		fmt.Println(err)
	}
}

// 销毁链接之后执行的钩子函数
func DoConnEnd(conn ziface.IConnection) {
	fmt.Println("===>DoConnEnd is called...")
	fmt.Println(">>Hook: conn ID = ", conn.GetConnID(), "has leaved.")
}

func main() {
	// 1.创建一个句柄(使用zinx的api)
	s := znet.NewServer("[zinx V0.9]")
	// 2.注册链接hook函数
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnEnd)
	// 3.给当前zinx框架添加router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	// 4.启动server
	s.Serve()
}