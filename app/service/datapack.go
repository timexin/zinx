package service

import (
	"bytes"
	"encoding/binary"
	"zinx/app/ifce"
)

type DataPack struct {
}

func NewDataPack() ifce.IDataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	// Datalen 4字节 +ID 4字节
	return 8
}

func (d *DataPack) Pack(msg ifce.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	// 讲datalen 写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//讲MsgId  写进databuff 中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 讲data数据写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (d *DataPack) Unpack(binaryData []byte) (ifce.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	// 只解压head信息，得到datalen，和msgID
	msg := &Message{}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	// 判断datalen是否已经超出定义的最大包长度
	//if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
	//	return nil, errors.New("too large msg data ")
	//}

	return msg, nil
}
