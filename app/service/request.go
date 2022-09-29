package service

import "zinx/app/ifce"

type Request struct {
	conn ifce.IConnection
	msg  ifce.IMessage
}

func NewRequest() ifce.IRequest {
	return &Request{}
}
func (r *Request) GetConnection() ifce.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
