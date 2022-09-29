package service

import (
	"fmt"
	"net"
	"zinx/app/ifce"
	"zinx/app/utils"
)

type zix struct {
	Name    string
	Version string
	Ip      string
	Port    string
	Router  ifce.IRouter
}

func (z *zix) Start() {

	go func() {
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
			conn := NewConn(acceptTCP, cid, z.Router)

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
	//TODO implement me
	panic("implement me")
}

func (z *zix) Serve() {
	z.Start()
	select {}
}

func (z *zix) AddRouter(router ifce.IRouter) {
	z.Router = router
	fmt.Println("Add Router Success")
}

func NewSvr(name string) ifce.Isv {
	return &zix{
		Name:    utils.GlobalObject.Name,
		Version: "tcp4",
		Ip:      utils.GlobalObject.Host,
		Port:    utils.GlobalObject.TcpPort,
	}
}
