package service

import (
	"fmt"
	"net"
	"zinx/app/ifce"
	"zinx/app/utils"
)

type zix struct {
	Name        string
	Version     string
	Ip          string
	Port        string
	MsgHandle   ifce.IMsgHandle
	ConnMgr     ifce.IConnManager
	OnConnStart func(conn ifce.IConnection)
	OnConnStop  func(conn ifce.IConnection)
}

func (z *zix) Start() {

	go func() {

		// 开启消息队列及woeker工作池
		z.MsgHandle.StartWorkerPool()
		// 获取tcp的attr
		fmt.Println("[start] server listen")
		fmt.Printf("[Zinx] Server Name :%s,listenner at IP :%s,Port: %s is starting", z.Name, z.Ip, z.Port)
		addr, err := net.ResolveTCPAddr(z.Version, fmt.Sprintf("%s:%s", z.Ip, z.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
		}
		// 监听
		tcp, err := net.ListenTCP(z.Version, addr)
		if err != nil {
			fmt.Println("listent: ", z.Version, " err", err)
		}
		var cid uint32
		cid = 0
		for {
			acceptTCP, err := tcp.AcceptTCP()
			if err != nil {
				fmt.Println("accept err", err)
				continue
			}
			fmt.Println(z.GetConnMgr().Len())
			fmt.Println(utils.GlobalObject.MaxConn)
			// 设置最大链接个数的判断
			if z.GetConnMgr().Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("Too Many Connections MaxConn =", utils.GlobalObject.MaxConn)
				acceptTCP.Close()
				continue
			}
			conn := NewConn(z, acceptTCP, cid, z.MsgHandle)

			cid++
			go conn.Start()

			//// 回现
			//go func() {
			//	for {
			//		buf := make([]byte, 512)
			//		read, err := acceptTCP.Read(buf)
			//		if err != nil {
			//			fmt.Println("recv buf err", err)
			//			continue
			//		}
			//		if _, err := acceptTCP.Write(buf[:read]); err != nil {
			//			fmt.Println("write back err ", err)
			//			continue
			//		}
			//	}
			//}()
		}

	}()

}

func (z *zix) Stop() {
	fmt.Println("[stop] ....")
	z.ConnMgr.ClearAll()
}

func (z *zix) Serve() {
	z.Start()
	select {}
}

func (z *zix) AddRouter(msgId uint32, router ifce.IRouter) {
	z.MsgHandle.AddRouter(msgId, router)
	fmt.Println("Add Router Success")
}

func (z *zix) GetConnMgr() ifce.IConnManager {
	return z.ConnMgr
}

func NewSvr(name string) ifce.Isv {
	return &zix{
		Name:      utils.GlobalObject.Name,
		Version:   "tcp4",
		Ip:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandle: NewMsgHandle(),
		ConnMgr:   NewConnManager(),
	}
}

func (z *zix) SetOnConnStart(f func(conn ifce.IConnection)) {
	z.OnConnStart = f
}

func (z *zix) SetOnConnStop(f func(conn ifce.IConnection)) {
	z.OnConnStop = f
}

func (z *zix) CallOnConnStart(conn ifce.IConnection) {
	if z.OnConnStart != nil {
		fmt.Println("---> Call onConnStart..")
		z.OnConnStart(conn)
	}
}

func (z *zix) CallOnConnStop(conn ifce.IConnection) {
	if z.OnConnStop != nil {
		fmt.Println("---> Call onConnStop..")
		z.OnConnStop(conn)
	}
}
