package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx_framework_demo/zinx/ziface"
)

/*
	链接管理模块
 */

type ConnManager struct {
	connections map[uint32] ziface.IConnection	//管理的链接集合
	connLock sync.RWMutex	//保护链接集合的读写锁
}

// 创建当前链接
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32] ziface.IConnection),
	}
}

// 添加链接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 添加共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将conn加入到ConnManager中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connID=", conn.GetConnID(),
		"connection add to ConnManager success! current conn num=", cm.Len())
}
// 删除链接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除链接信息
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connID=", conn.GetConnID(),
		"connection has been removed from ConnManager success! current conn num=", cm.Len())
}
// 根据connID获取链接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源，加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		// 找到了对应的conn
		return conn, nil
	} else {
		return nil, errors.New("connection NOT FOUND!")
	}
}
// 得到当前链接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}
// 停止并删除所有连接
func (cm *ConnManager) ClearConn() {
	// 保护共享资源，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除conn并终止conn的工作
	for connID, conn := range cm.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.connections, connID)
	}
	fmt.Println("clear all conncetons success! current conn num=", cm.Len())
}