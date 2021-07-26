package main

import (
	"zinx_framework_demo/zinx/znet"
)

/*
	基于zinx来开发的服务器端应用程序
 */

func main() {
	// 1.创建一个句柄(使用zinx的api)
	s := znet.NewServer("[zinx V0.1]")
	// 2.启动server
	s.Serve()
}