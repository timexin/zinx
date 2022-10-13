package service

import (
	"fmt"
	"zinx/app/ifce"
	"zinx/app/utils"
)

type MsgHandle struct {
	Apis map[uint32]ifce.IRouter
	//负责worker取任务的消息队列
	TaskQueue []chan ifce.IRequest
	//worker 数量
	WorkerPoolSize uint32
}

func NewMsgHandle() ifce.IMsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ifce.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ifce.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (m *MsgHandle) DoMsgHandler(request ifce.IRequest) {
	router, ok := m.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgId(), "router is null")
	}
	router.Handle(request)
}

func (m *MsgHandle) AddRouter(msgID uint32, router ifce.IRouter) {
	// 1.判断当前msg绑定的API处理方法是否存在
	if _, ok := m.Apis[msgID]; ok {

	} else {
		m.Apis[msgID] = router
	}
}

func (m *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan ifce.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前的worker ,足赛等待从cnanel传递进来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// 启动一个Worker 工作流程
func (m *MsgHandle) StartOneWorker(workerID int, taskQueue chan ifce.IRequest) {
	fmt.Println("Worker ID = ", workerID, "is started ...")
	for {
		select {
		//如果有消息过来，就是一个queset
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQuque 将消息平均分配给不同的worker
func (m *MsgHandle) SendMsgToTaskQuque(request ifce.IRequest) {
	// 根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), "request MsgID =", request.GetMsgId(), "to WorkID =", workerID)
	m.TaskQueue[workerID] <- request
}
