package service

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 只是负责测试datapack 拆包 封包的方法
func TestDataPack(t *testing.T) {
	// 模拟的服务器
	//  1.创建socketTcp
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}
	go func() {
		//  2.从客户端读取数据，拆包处理
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept err:", err.Error())
			}
			go func(conn net.Conn) {
				//处理客户端的请求
				// -----》拆包的过程
				dp := NewDataPack()
				for {
					// 第一次从conn读，把包的header读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err")
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//msg 是有数据的 进行第二次读取
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack err", err)
						}
						fmt.Println("-------> RecvMsgId:", msg.Id, ",datalen =", msg.DataLen, ",data = ", msg.Data)
					}
				}
			}(conn)
		}
	}()
	/*
		模拟客户端
	*/
	dial, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	dps := NewDataPack()
	// 模拟粘包过程，封装2个msg一起发送
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	pack1, err := dps.Pack(msg1)
	if err != nil {
		fmt.Println("cient pack msg1 err ", err)
	}

	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'z', 'i', 'n', 'x', 'h', 'e', 'l'},
	}
	pack2, err := dps.Pack(msg2)
	if err != nil {
		fmt.Println("cient pack msg1 err ", err)
	}
	pack1 = append(pack1, pack2...)
	dial.Write(pack1)
	select {}
}
