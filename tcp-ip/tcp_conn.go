package main

import (
	"fmt"
	"log"
	"tcpip/tcpip"
	"time"
)

func main() {
	dest := "192.168.1.1"
	var port uint16 = 49152

	// TCP接続を開始するためのSYNパケットを作成
	syn := tcpip.TCPIP{
		DestIP:   dest,
		DestPort: port,
		TcpFlag:  "SYN",
	}

	// 送信元IPアドレスとポート番号を設定
	ifaceName := "en0"
	sendfd := tcpip.NewTCP(ifaceName, port)
	defer sendfd.Close()

	// 3Way HandshakeでTCP接続を確立
	ack, err := sendfd.StartTCPConnection(syn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("TCP Connection is success!!\n\n")
	time.Sleep(10 * time.Millisecond)

	fin := tcpip.TCPIP{
		DestIP:    dest,
		DestPort:  port,
		TcpFlag:   "FINACK",
		SeqNumber: ack.SeqNumber,
		AckNumber: ack.AckNumber,
	}
	_, err = sendfd.StartTCPConnection(fin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("TCP Connection Close is success!!\n")
}
