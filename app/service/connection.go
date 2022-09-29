package service

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/app/ifce"
)

type Connection struct {
	// 当前链接的socket TCP 套接字
	Conn *net.TCPConn
	// 链接ID
	ConnID uint32
	// 当前链接的状态
	IsClosed bool
	// 处理业务方法API
	Router ifce.IRouter
	// 告知当前链接已经退出的channel
	ExitChan chan bool
}

func NewConn(conn *net.TCPConn, connID uint32, router ifce.IRouter) ifce.IConnection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		IsClosed: false,
		ExitChan: make(chan bool, 1),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running ...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit,remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读数据
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}
		//创建一个拆包的对象
		dp := NewDataPack()
		//读取客户端的Msg Head 二进制流
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg header err", err)
			break
		}
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetMsgData(data)
		req := &Request{
			conn: c,
			msg:  msg,
		}
		go func(request ifce.IRequest) {
			//		c.Router.PreHandle(request)
			c.Router.Handle(request)
			//		c.Router.PostHandle(request)
		}(req)
		// 从路由中，找到该链接的Conn对应的router调用

	}
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("Connection is closed")
	}
	dp := NewDataPack()
	msg := new(Message)
	msg.Id = 1
	msg.Data = data
	msg.DataLen = uint32(len(data))
	pack, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("Pack error msg id =", msg.Id)
		return err
	}
	if _, err := c.Conn.Write(pack); err != nil {
		return err
	}
	return nil
}

func (c *Connection) Start() {
	fmt.Println("Conn start() ... ConnId = ", c.ConnID)
	go c.StartReader()
}

func (c *Connection) Stop() {
	if c.IsClosed {
		return
	}
	c.IsClosed = true
	c.Conn.Close()
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
