package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx_framework_demo/zinx/ziface"
)

/*
	存储一切有关Zinx框架的全局参数，供其它模块使用。
	一些参数是可以通过zinx.json由用户进行配置
*/
type GlobalObj struct {
	/*	Server	*/
	TcpServer ziface.IServer // 当前Zinx全局的Server对象
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            // 当前服务器主机监听的端口号
	Name      string         // 当前服务器的名称
	/*	Zinx	*/
	Version          string // 当前Zinx的版本号
	MaxConn          int    // 当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 // 当前Zinx框架数据包的最大值
	WorkerPoolSize   uint32 // 当前业务工作Worker池的Goroutine数量
	MaxWorkerTaskLen uint32 // Zinx框架允许用户最大开辟Worker数
}

/*
	定义一个全局的对外Globalobj对象
*/
var GlobalObject *GlobalObj

// 从配置文件去加载用于自定义的参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json文件数据解析到GlobalObject
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供一个init方法，用来初始化当前的GlobalObject全局对象
func init() {
	// 如果配置文件没有加载，就是默认的值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.8",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	// 应该尝试从配置文件中去加载默认的值
	GlobalObject.Reload()
}
