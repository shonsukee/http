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
}
