package main

import "tcpip/tcpip"

func main() {
	// TODO: 宛先IPアドレスを変更して試すとより良い！
	dstIP := "192.168.1.1"

	// TODO: `networksetup -listallhardwareports`コマンドで確認
	interfaceName := "en0"

	// ARPリクエストを送信
	if err := tcpip.ArpRequest(interfaceName, dstIP); err != nil {
		panic(err)
	}

	// UDPパケットを送信(例:DNSポート53にメッセージを送信)
	srcPort := uint16(49152)
	dstPort := uint16(53)
	message := "Hello UDP!"

	if err := tcpip.UdpRequest(interfaceName, dstIP, srcPort, dstPort, message); err != nil {
		panic(err)
	}
}
