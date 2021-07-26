package znet

import (
	"fmt"
	"strconv"
	"zinx_framework_demo/zinx/utils"
	"zinx_framework_demo/zinx/ziface"
)

/*
	消息处理模块的实现
*/

type MsgHandle struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责worker读取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作worker池的worker数量
	WorkerPoolSize uint32
}

// 创建MsgHandle的方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 从全局配置中获取
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1.从Request找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("API msgID=", request.GetMsgID(), " not found!")
	}
	// 2.从对应的MsgID调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1.判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		//id已经注册
		panic("repeat API, msgID=" + strconv.Itoa(int(msgID)))
	}

	// 2.添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID=", msgID, " succsess!")
}

// 启动一个Worker工作池(开启工作池的动作只能发生一次，一个Zinx框架只有一个worker工作池)
func (mh *MsgHandle) StartWorkerPool() {
	// 根据workerPoolSize分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		// 1.给当前的worker对应的channel消息队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2.启动当前的worker，阻塞等待消息从channel中到来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerID=", workerID, " is starting...")
	// 不断地则色等待消息队列的消息
	for {
		select {
			// 如果有消息过来
			case request := <-taskQueue:
				mh.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueue，由worker来进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1.将消息平均分配给不同的worker
	// 根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request MsgID=",
		request.GetMsgID(), " to WorkerID=", workerID)

	// 2.将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}