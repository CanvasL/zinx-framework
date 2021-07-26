package main

import (
	"fmt"
	"net"
	"time"
)

/*
	模拟客户端
 */
func main() {
	fmt.Println("Client Start...")

	time.Sleep(time.Second * 1)
	// 1.直接链接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start failed, err:", err)
	}
	// 2.链接调用Write方法，写数据
	for {
		_, err := conn.Write([]byte("Hello Zinx V0.1!"))
		if err != nil {
			fmt.Println("write conn failed, err:", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf failed, err:", err)
			return
		}

		fmt.Printf("server call back: %s, cnt = %d\n", buf[:cnt], cnt)

		// cpu阻塞
		time.Sleep(time.Second * 1)
	}
}