package service

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/app/ifce"
	"zinx/app/utils"
)

type Connection struct {
	//当前Conn属于哪个Server
	TcpServer ifce.Isv
	// 当前链接的socket TCP 套接字
	Conn *net.TCPConn
	// 链接ID
	ConnID uint32
	// 当前链接的状态
	IsClosed bool
	// 处理业务方法API
	MsgHandler ifce.IMsgHandle
	// 告知当前链接已经退出的channel
	ExitChan chan bool
	// 读写chan
	MsgChan      chan []byte
	ProPerty     map[string]interface{}
	ProPertyLock sync.RWMutex
}

func NewConn(server ifce.Isv, conn *net.TCPConn, connID uint32, MsgHandler ifce.IMsgHandle) ifce.IConnection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: MsgHandler,
		IsClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgChan:    make(chan []byte),
		TcpServer:  server,
		ProPerty:   map[string]interface{}{},
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
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
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQuque(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
			// 从路由中，找到该链接的Conn对应的router调用
		}

	}
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("Connection is closed")
	}
	dp := NewDataPack()
	msg := new(Message)
	msg.Id = msgId
	msg.Data = data
	msg.DataLen = uint32(len(data))
	pack, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("Pack error msg id =", msg.Id)
		return err
	}
	//if _, err := c.Conn.Write(pack); err != nil {
	//	return err
	//}
	c.MsgChan <- pack
	return nil
}

// StartWriter 专门发送给客户端的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Gortine is running")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit]")
	for {
		select {
		case data := <-c.MsgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error", err)
				return
			}
		case <-c.ExitChan:
			// 代表Reader已经退出，此时Writer也要退出
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn start() ... ConnId = ", c.ConnID)
	go c.StartReader()
	go c.StartWriter()
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	if c.IsClosed {
		return
	}
	c.IsClosed = true
	c.TcpServer.CallOnConnStop(c)
	c.Conn.Close()
	c.ExitChan <- true
	c.TcpServer.GetConnMgr().Remove(c)
	close(c.ExitChan)
	close(c.MsgChan)
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

//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.ProPertyLock.Lock()
	defer c.ProPertyLock.Unlock()
	c.ProPerty[key] = value
}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.ProPertyLock.RLock()
	defer c.ProPertyLock.RUnlock()
	if value, ok := c.ProPerty[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no protery found")
	}

}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.ProPertyLock.Lock()
	defer c.ProPertyLock.Unlock()
	delete(c.ProPerty, key)

}
