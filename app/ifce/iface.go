package ifce

type Isv interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgId uint32, router IRouter)
	GetConnMgr() IConnManager
	SetOnConnStart(func(conn IConnection))
	SetOnConnStop(func(conn IConnection))
	CallOnConnStart(conn IConnection)
	CallOnConnStop(conn IConnection)
}
