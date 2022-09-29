package ifce

type IRequest interface {
	GetConnection() IConnection

	GetData() []byte

	GetMsgId() uint32
}
