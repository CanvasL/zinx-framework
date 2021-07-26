 package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx_framework_demo/zinx/znet"
)

/*
	模拟客户端
 */
func main() {
	fmt.Println("Client1 Start...")

	time.Sleep(time.Second * 1)
	// 1.直接链接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start failed, err:", err)
	}
	// 2.链接调用Write方法，写数据
	for {
		// 发送封包的message信息, MsgID:0
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(1, []byte("ZinxV0.6 client hello test")))
		if err != nil {
			fmt.Println("Pack msg failed, ", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("conn Write failed, ", err)
			return
		}

		// 服务器就应该给我回复一个message数据， MsgID:1 ping...ping...ping...
		// 1.先读取流中的header部分，得到ID和dataLen
		binaryHeader := make([]byte, dp.GetHeaderLen())
		if _, err := io.ReadFull(conn, binaryHeader); err != nil {
			fmt.Println("read header failed, ", err)
			break
		}
		//将二进制的header拆包到msg结构体中
		msgHeader, _ := dp.Unpack(binaryHeader)
		if msgHeader.GetMsgLen() > 0 {
			//说明data中有数据
			// 2.再根据DataLen进行第二次读取，将data读取出来
			msg := msgHeader.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data failed, ", err)
				return
			}

			fmt.Println(">>>Recv Server Msg, MsgId=", msg.Id,
				", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		// cpu阻塞
		time.Sleep(time.Second * 3)
	}
}