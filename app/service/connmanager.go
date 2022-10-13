package service

import (
	"errors"
	"fmt"
	"sync"
	"zinx/app/ifce"
)

type ConnManager struct {
	connections map[uint32]ifce.IConnection //管理的链接集合
	connLock    sync.RWMutex                //保护连接集合的读写锁
}

func NewConnManager() ifce.IConnManager {
	return &ConnManager{
		connections: make(map[uint32]ifce.IConnection),
	}
}

func (c *ConnManager) Add(conn ifce.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	// 将conn加入到ConnManager中
	c.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager successfully: conn num =", c.Len())
}

func (c *ConnManager) Remove(conn ifce.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	delete(c.connections, conn.GetConnID())
	fmt.Println("connID = ", conn.GetConnID(), "remove from ConnManager successfully ")
}

func (c *ConnManager) Get(connID uint32) (ifce.IConnection, error) {
	// 加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}

}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearAll() {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	//删除conn,并请知conn的工作
	for connId, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connId)
	}
	fmt.Println("Clear All connetions ")
}
