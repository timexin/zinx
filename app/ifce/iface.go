package ifce

type Isv interface {
	Start()
	Stop()
	Serve()
	AddRouter(router IRouter)
}
