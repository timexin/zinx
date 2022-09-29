package ifce

type IMessage interface {
	GetMsgId() uint32
	GetMsgLen() uint32
	GetData() []byte
	SetMsgId(uint32)
	SetMsgLen(uint32)
	SetMsgData([]byte)
}
