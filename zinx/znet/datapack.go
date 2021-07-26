package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx_framework_demo/zinx/utils"
	"zinx_framework_demo/zinx/ziface"
)

// 封包、拆包的具体模块
type DataPack struct{}

// 拆包封包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包的header长度方法
func (dp *DataPack) GetHeaderLen() uint32 {
	// DataLen uint32(4 Byte) + ID uint32(4 Byte) = 8 Byte
	return 8
}

// 封包方法
//		写入格式：|dataLen|Id|data|
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放byte字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将msgId写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 将data数据写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法
// 先读取固定长度的header，然后根据header中的dataLen读取一定字节的消息内容
func (dp *DataPack) Unpack(binaryData[]byte) (ziface.IMessage, error) {
	// 创建一个读取binary数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	// 只解压header信息，得到dataLen和msgId
	msg := &Message{}
	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	// 判断dataLen是否已经超过我们允许的最大包长度
	if (utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize) {
		return nil, errors.New("recv msg data too long!")
	}

	return msg, nil
}
