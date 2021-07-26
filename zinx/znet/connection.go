package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx_framework_demo/zinx/utils"
	"zinx_framework_demo/zinx/ziface"
)

/*
	链接模块
*/
type Connection struct {
	// 当前conn隶属于哪个Server
	TcpServer ziface.IServer
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn
	// 链接的ID
	ConnID uint32
	// 当前的链接状态protobuf
	isClosed bool
	// 告知当前链接已经退出/停止的channel,由reader告知writer
	ExitChan chan bool
	// 无缓冲的通道，用于读、写goroutine通信的消息
	msgChan chan []byte
	// 消息的管理MsgID和对应的业务处理api
	MsgHandler ziface.IMsgHandle
	// 链接属性集合
	property map[string]interface{}
	// 保护链接属性的锁
	propertyLock sync.RWMutex
}

// 初始化链接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		property: make(map[string]interface{}),
		MsgHandler: msgHandler,
	}

	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine] is running...")
	defer fmt.Println("connID=", c.ConnID, ", reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 创建一个拆包解包对象
		dp := NewDataPack()
		// 读取客户端的MsgHeader二进制流(8个字节)
		headData := make([]byte, dp.GetHeaderLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg header failed, ", err)
			break
		}
		// 拆包，得到msgID和msgDatalen，放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack failed, ", err)
			break
		}
		// 根据dataLen再次读取Data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data failed, ", err)
				break
			}
		}
		msg.SetMsgData(data)
		// 得到当前conn的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 根据绑定好的MsgID找到对应的业务处理API，执行
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 写消息的Goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine] is running! ")
	defer fmt.Println("connID=", c.ConnID, ", writer is exit, remote addr is ", c.RemoteAddr().String())

	// 不断地等待channel的消息，进行回写
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data failed, ", err)
				return
			}
		case <-c.ExitChan:
			// 代表reader已经退出，此时writer也要退出
			return
		}
	}
}

// 启动链接
func (c *Connection) Start() {
	fmt.Println("connection start, ConnID=", c.ConnID)

	// 启动从当前链接的读数据的业务
	go c.StartReader()
	// 启动从当前链接写数据的业务
	go c.StartWriter()

	// 调用开发者注册的hook函数
	c.TcpServer.CallOnConnStart(c)
}

// 停止链接
func (c *Connection) Stop() {
	fmt.Println("connection stop, ConnID=", c.ConnID)

	// 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 调用开发者注册的hook函数
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket链接
	c.Conn.Close()

	// 告知writer关闭
	c.ExitChan <- true

	// 将当前链接从ConnMgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前链接模块的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接模块的链接id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的 TCP状态 ip port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据，将数据先封包，再发送给远程的客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection is closed when sending msg")
	}

	// 将data进行封包,顺序：MsgDataLen|MsgID|Data
	dp := NewDataPack()
	msg := NewMsgPackage(msgId, data)

	// binaryMsg的格式就是：MsgDataLen|MsgID|Data的二进制流
	binaryMsg, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("pack msg failed, msgID = ", msgId)
		return errors.New("pack msg failed")
	}

	// 将数据发送给管道，让writer得到消息并写回客户端
	c.msgChan <- binaryMsg

	return nil
}

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 添加一个链接属性
	c.property[key] = value
}

// 获取链接属性
func (c *Connection) GetProperty(key string)(interface{}, error) {
	c.propertyLock.RLock()
	defer c. propertyLock.RUnlock()

	// 读取属性
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("Can not find property.")
	}
}

// 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 删除属性
	delete(c.property, key)
}