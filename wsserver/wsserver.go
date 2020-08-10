package main

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/chengchaos/go-learning/helper"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}

	for {
		func() {
			conn, err := ln.Accept()
			if err != nil {
				log.Println("Accept get error :", err)
				return
			}

			for {
				handleConnection(conn)
			}
		}()
	}

	fmt.Println(helper.ChineseGBK("启动了"))

}

func handleConnection(conn net.Conn) {
	content := make([]byte, 1024)

	_, err := conn.Read(content)

	if err != nil {
		log.Println("Read content get error:", err)
	}

	log.Printf("content => '%v'\n", string(content))
	isHttp := false

	if string(content[0:3]) == "GET" {
		isHttp = true
	}

	if isHttp {
		headers := parseHandshake(string(content))
		log.Println("headers =>", headers)

		secWebsocketKey := headers["Sec-WebSocket-Key"]

		//
		guid := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

		// 计算 Sec-WebSocket-Accept
		h := sha1.New()
		acceptRaw := secWebsocketKey + guid
		log.Println("accept raw =>", acceptRaw)
		io.WriteString(h, acceptRaw)

		accept := make([]byte, 28)
		base64.StdEncoding.Encode(accept, h.Sum(nil))
		log.Println(string(accept))

		responseTemp := `
HTTP/1.1 101 Switching Protocols
Sec-WebSocket-Accept: %s
Connection: Upgrade
Upgrade: websocket

`
		response := fmt.Sprintf(responseTemp, string(accept))
		log.Println("response =>", response)

		if length, err := conn.Write([]byte(response)); err != nil {
			log.Println("Write but get an error", err)
		} else {
			log.Println("send len =>", length)
		}

		wssocket := NewWsSocket(conn)

		for {
			data, err := wssocket.ReadIframe()
			if err != nil {
				log.Println("ReadIframe error =>", err)
			}
			log.Println(helper.ChineseGBK("接收 data =>", string(data)))

			err = wssocket.SendIframe([]byte("good"))
			if err != nil {
				log.Println("Sendiframe error =>", err)
			}
			log.Println(helper.ChineseGBK("回复 data OK!"))
		}
	} else {
		log.Println(string(content))
	}
}

type WsSocket struct {
	MaskingKey []byte
	Conn       net.Conn
}

func NewWsSocket(conn net.Conn) *WsSocket {
	return &WsSocket{Conn: conn}
}

func (ws *WsSocket) SendIframe(data []byte) error {


	dataLength := len(data)

	if dataLength >= 125 {
		return errors.New("我们只处理 data ,< 125 的(⊙﹏⊙)")
	}

	maskedData := make([]byte, dataLength)

	for i := 0; i < dataLength; i++ {
		if ws.MaskingKey != nil {
			maskedData[i] = data[i] ^ ws.MaskingKey[i % 4]
		} else {
			maskedData[i] = data[i]
		}
	}

	ws.Conn.Write([]byte{0x81})

	var payLenByte byte
	if ws.MaskingKey != nil && len(ws.MaskingKey) != 4 {
		payLenByte = byte(0x80) | byte(dataLength)
		ws.Conn.Write([]byte{payLenByte})
		ws.Conn.Write(ws.MaskingKey)
	} else {
		payLenByte = byte(0x00) | byte(dataLength)
		ws.Conn.Write([]byte{payLenByte})
	}

	_, err := ws.Conn.Write(data)
	return err
}

/*-----------为了便于理解，在这里吧数据帧格式粘出来-------------------
0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-------+-+-------------+-------------------------------+
|F|R|R|R| opcode|M| Payload len |    Extended payload length    |
|I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
|N|V|V|V|       |S|             |   (if payload len==126/127)   |
| |1|2|3|       |K|             |                               |
+-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
|     Extended payload length continued, if payload len == 127  |
+ - - - - - - - - - - - - - - - +-------------------------------+
|                               |Masking-key, if MASK set to 1  |
+-------------------------------+-------------------------------+
| Masking-key (continued)       |          Payload Data         |
+-------------------------------- - - - - - - - - - - - - - - - +
:                     Payload Data continued ...                :
+ - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
|                     Payload Data continued ...                |
+---------------------------------------------------------------+
--------------------------------------------------------------------*/
func (ws *WsSocket) ReadIframe() (data []byte, err error) {

	// 第一个字节 ：
	// FIN + RSV1-3 + OPCODE
	opcodeByte := make([]byte, 1)
	_, err = ws.Conn.Read(opcodeByte)

	if err != nil {
		return
	}

	fin := opcodeByte[0] >> 7
	rsv1 := opcodeByte[0] >> 6 & 1
	rsv2 := opcodeByte[0] >> 5 & 1
	rsv3 := opcodeByte[0] >> 4 & 1
	opcode := opcodeByte[0] & 15

	log.Printf("fin => %v, rsv => %v, %v, %v, opcode => %v\n", fin, rsv1, rsv2, rsv3, opcode)

	payloadLenByte := make([]byte, 1)
	_, err = ws.Conn.Read(payloadLenByte)
	if err != nil {
		return
	}
	payloadLen := int(payloadLenByte[0] & 0x7F)
	mask := payloadLenByte[0] >> 7

	log.Printf("payloadLen => %d, mask => %v\n", payloadLen, mask)

	if payloadLen == 127 {
		extendedByte := make([]byte, 8)
		_, err = ws.Conn.Read(extendedByte)
		if err != nil {
			return
		}
	}

	maskingByte := make([]byte, 4)
	if mask == 1 {
		_, err = ws.Conn.Read(maskingByte)
		if err != nil {
			return
		}
		ws.MaskingKey = maskingByte
	}

	payloadDataByte := make([]byte, payloadLen)
	_, err = ws.Conn.Read(payloadLenByte)
	if err != nil {
		return
	}
	log.Println("data =>", payloadDataByte)

	dataByte := make([]byte, payloadLen)
	for i := 0; i < payloadLen; i++ {
		if mask == 1 {
			dataByte[i] = payloadDataByte[i] ^ maskingByte[i % 4]
		} else {
			dataByte[i] = payloadDataByte[i]
		}
	}

	log.Printf("fin == 1 => %v\n", fin == 1)

	if fin == 1 {
		data = dataByte
		return
	}

	nextData, err := ws.ReadIframe()
	if err != nil {
		return
	}

	data = append(data, nextData...)
	return
}

func parseHandshake(content string) map[string]string {

	headers := make(map[string]string, 10)
	lines := strings.Split(content, "\r\n")

	for _, line := range lines {
		if len(line) > 0 {
			words := strings.Split(line, ":")
			if len(words) == 2 {
				headers[strings.Trim(words[0], " ")] = strings.Trim(words[1], " ")
			}
		}
	}

	return headers
}