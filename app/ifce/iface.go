package ifce

type Isv interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgId uint32, router IRouter)
}
