package main

import (
	"fmt"
	"zinx/app/ifce"
	"zinx/app/service"
)

func main() {
	svr := service.NewSvr("demo")
	// 添加自定义模版
	svr.AddRouter(&PingRouter{})
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

func (this *PingRouter) PreHandle(request ifce.IRequest) {

}

func (this *PingRouter) PostHandle(request ifce.IRequest) {

}