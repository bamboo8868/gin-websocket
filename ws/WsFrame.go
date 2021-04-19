package ws

type WsFrame struct {
	Fin        byte
	Rsv1       byte
	Rsv2       byte
	Rsv3       byte
	Opcode     byte
	Mask       byte
	PayloadLen byte
	MaskingKey []byte
	Data       []byte
}
