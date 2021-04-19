package ws

import "net"

type WsConn struct {
	Conn net.Conn
}

func Init(c net.Conn) *WsConn {
	conn := &WsConn{}
	conn.Conn = c
	return conn
}

func (conn *WsConn) ReadMessage() *WsFrame{
	frame := new(WsFrame)
	opByte := make([]byte, 1)
	conn.Conn.Read(opByte)

	fin := opByte[0] >> 7 & 1
	rsv1 := opByte[0] >> 6 & 1
	rsv2 := opByte[0] >> 5 & 1
	rsv3 := opByte[0] >> 4 & 1

	opcode := opByte[0] & 15

	frame.Fin = fin
	frame.Rsv1 = rsv1
	frame.Rsv2 = rsv2
	frame.Rsv3 = rsv3
	frame.Opcode = opcode

	//byte2 包含 mask payload len
	byte2 := make([]byte, 1)
	conn.Conn.Read(byte2)

	mask := byte2[0] >> 7 & 1

	len := byte2[0] & 0b01111111

	frame.Mask = mask
	frame.PayloadLen = len

	//只处理小于126字节情况


	//读取mask key
	if mask == 1 {
		//有掩码 数据需要掩码处理
		maskKey := make([]byte, 4)

		conn.Conn.Read(maskKey)
		frame.MaskingKey = maskKey
		data := make([]byte, len)
		conn.Conn.Read(data)
		frame.Data = data

		for i := 0; i < int(len); i++ {
			data[i] = data[i] ^ maskKey[i%4]
		}
		frame.Data = data
	} else {
		//无掩码 直接读取数据
		data := make([]byte, len)

		conn.Conn.Read(data)
		frame.Data = data
	}

	return frame

}

func (conn *WsConn) SendMessage(msg []byte) {

}
