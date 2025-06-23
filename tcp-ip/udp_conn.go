package main

import "tcpip/tcpip"

func main() {
	interfaceName := "en0"
	dstIP := "192.168.1.1"

	// UDPパケットを送信(例:DNSポート53にメッセージを送信)
	srcPort := uint16(49152)
	dstPort := uint16(53)
	message := "Hello UDP!"

	if err := tcpip.UdpRequest(interfaceName, dstIP, srcPort, dstPort, message); err != nil {
		panic(err)
	}
}
