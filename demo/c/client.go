package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/app/service"
)

func main() {
	fmt.Println("client start")
	time.Sleep(time.Second)
	dial, err := net.Dial("tcp", "127.0.0.1:7766")
	if err != nil {
		fmt.Println("conn err ", err)
	}
	for {
		dp := service.NewDataPack()
		msg := new(service.Message)
		msg.Id = 0
		msg.Data = []byte("Zinx client test message")
		msg.DataLen = uint32(len(msg.Data))
		pack, err := dp.Pack(msg)
		if err != nil {
			fmt.Println("pack err", err)
		}
		if _, err := dial.Write(pack); err != nil {
			fmt.Println("write err")
		}
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(dial, headData); err != nil {
			fmt.Println("read head err", err)
		}
		unpack, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("client ubpack msgHead err", err)
		}
		if unpack.GetMsgLen() > 0 {
			realMsg := unpack.(*service.Message)
			realMsg.Data = make([]byte, realMsg.GetMsgLen())
			if _, err := io.ReadFull(dial, realMsg.Data); err != nil {
				fmt.Println("read msg data error ", err)
				return
			}
			fmt.Println("---->Recv Server id = ", realMsg.Id, ",len =", realMsg.DataLen, ",data = ", string(realMsg.Data))
		}

		time.Sleep(time.Second)
	}

}
