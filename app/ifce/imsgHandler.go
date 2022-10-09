package ifce

type IMsgHandle interface {
	// 执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)

	//为消息添加具体的处理方法
	AddRouter(msgID uint32, router IRouter)
}
