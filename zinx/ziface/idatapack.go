package ziface

/*
	封包、拆包模块
	直接面向TCP连接中的数据流，用于处理TCP粘包问题
*/

type IDataPack interface {
	// 获取包的header长度方法
	GetHeaderLen() uint32
	// 封包方法
	Pack(msg IMessage)([]byte, error)
	// 拆包方法
	Unpack([]byte)(IMessage, error)
}