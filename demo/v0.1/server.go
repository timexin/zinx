package main

import (
	"fmt"
	"zinx/app/ifce"
	"zinx/app/service"
)

func main() {
	svr := service.NewSvr("demo")
	svr.SetOnConnStart(func(conn ifce.IConnection) {
		fmt.Println("----->DoConnectionBegin is Called ..")
		if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
			fmt.Println(err)
		}
	})
	svr.SetOnConnStop(func(conn ifce.IConnection) {
		fmt.Println("----->DoConnectionClosed is Called ..")
	})
	// 添加自定义模版
	svr.AddRouter(0, &PingRouter{})
	svr.AddRouter(1, &HelloRouter{})
	svr.Serve()
}

type PingRouter struct {
	ifce.IRouter
}

func (this *PingRouter) Handle(request ifce.IRequest) {
	fmt.Println("Call Router Handle...\n")
	fmt.Println("recv from client:", string(request.GetData()))
	err := request.GetConnection().SendMsg(request.GetMsgId(), []byte("ping..ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	ifce.IRouter
}

func (this *HelloRouter) Handle(request ifce.IRequest) {
	fmt.Println("Call Router Handle...\n")
	fmt.Println("recv from client:", string(request.GetData()))
	err := request.GetConnection().SendMsg(request.GetMsgId(), []byte("hello..hello...hello"))
	if err != nil {
		fmt.Println(err)
	}
}

func (this *PingRouter) PreHandle(request ifce.IRequest) {

}

func (this *PingRouter) PostHandle(request ifce.IRequest) {

}
