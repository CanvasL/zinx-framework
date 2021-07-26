package ziface

/*
	IRequest接口：
	实际上是把客户端请求的链接信息和请求的数据，包装到了一个Request中
 */

type IRequest interface {
	// 得到当前链接
	GetConnection() IConnection

	// 得到当前数据
	GetData() []byte

	// 得到当前请求消息的ID
	GetMsgID() uint32
}