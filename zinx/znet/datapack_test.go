package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 只是负责测试datapack拆包、封包的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟的服务器
	 */
	// 1.创建socketTCP
	listenner, err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("server listen failed, ", err)
	}

	// 创建一个go承载负责从客户端处理业务
	go func() {
		// 2.从客户端读取数据，拆包处理
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept failed, ", err)
			}
			go func(conn net.Conn) {
				// 处理客户端的请求
				//拆包的过程
				dp := NewDataPack()
				for {
					//1.先读包的header
					headData := make([]byte, dp.GetHeaderLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read header failed, ", err)
						break
					}
					msgHeader, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack failed, ", err)
						return
					}
					if msgHeader.GetMsgLen() > 0 {
						// msg是有数据的，需要第二次读取
						//2.根据header里的dataLen，读取data内容
						msg := msgHeader.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						// 根据dataLen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data failed, ", err)
							return
						}

						// 完整的一个消息已经读取完毕
						fmt.Println(">>>Recv MsgID: ", msg.Id, ", dataLen= ", msg.DataLen,
							", data= ", string(msg.Data))
					}
				}
			} (conn)
		}
	}()


	/*
		模拟客户端
	 */
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client dial failed, ", err)
		return
	}

	// 创建一个封包对象dp
	dp := NewDataPack()

	// 模拟粘包过程，封装两个msg一同发送
	// 封装第一个msg1包
	msg1 := &Message{
		Id:1,
		DataLen:4,
		Data: []byte{'z','i','n','x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 failed, ", err)
		return
	}
	// 封装第二个msg2包
	msg2 := &Message{
		Id:2,
		DataLen:7,
		Data: []byte{'h','l', 'l','z','i','n','x'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 failed, ", err)
		return
	}
	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)
	// 一起发送给server
	conn.Write(sendData1)

	// 客户端阻塞
	select{}
}