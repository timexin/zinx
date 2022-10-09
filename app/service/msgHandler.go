package service

import (
	"fmt"
	"zinx/app/ifce"
)

type MsgHandle struct {
	Apis map[uint32]ifce.IRouter
}

func NewMsgHandle() ifce.IMsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ifce.IRouter),
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
