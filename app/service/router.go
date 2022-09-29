package service

import "zinx/app/ifce"

type BaseRouter struct {
}

func NewBaseRouter() ifce.IRouter {
	return &BaseRouter{}
}

func (b *BaseRouter) PreHandle(request ifce.IRequest) {

}

func (b *BaseRouter) Handle(request ifce.IRequest) {

}

func (b *BaseRouter) PostHandle(request ifce.IRequest) {

}
