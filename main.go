package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"gin-websocket-demo/ws"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

const WSGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func IsWsConn(headers http.Header) bool {
	if _, ok := headers["Upgrade"]; !ok || headers["Upgrade"][0] != "websocket" {
		return false
	}

	if _, ok := headers["Sec-Websocket-Key"]; !ok {
		return false
	}

	if _, ok := headers["Sec-Websocket-Version"]; !ok || headers["Sec-Websocket-Version"][0] != "13" {
		return false
	}

	return true
}

func WsHandShake(ctx *gin.Context) {

	key := ctx.GetHeader("Sec-Websocket-Key")
	hash := sha1.New()
	hash.Write([]byte(key))
	hash.Write([]byte(WSGUID))
	secAccept := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	conn, buf, _ := ctx.Writer.Hijack()

	var buffer bytes.Buffer
	buffer.Write([]byte("HTTP/1.1 101 Switching Protocols\r\n"))
	buffer.Write([]byte("Upgrade: websocket\r\n"))
	buffer.Write([]byte("Connection: Upgrade\r\n"))
	buffer.Write([]byte("Sec-Websocket-Accept: " + secAccept + "\r\n\r\n"))
	conn.Write(buffer.Bytes())
	go handleWsConn(conn, buf)
}
func TlsHandShake() {

}

func handleWsConn(conn net.Conn, buf *bufio.ReadWriter) {
	wsConn := ws.Init(conn)
	for {
		frame := wsConn.ReadMessage()
		fmt.Println(string(frame.Data))
		wsConn.SendMessage()
	}

}



func main() {
	server := gin.Default()

	server.GET("/ws", func(ctx *gin.Context) {

		res := IsWsConn(ctx.Request.Header)

		if res {
			WsHandShake(ctx)
		} else {
			fmt.Println("this is http request")
			ctx.JSON(200, gin.H{"msg": "无效websocket请求"})
		}

	})

	server.Run(":8888")
}
