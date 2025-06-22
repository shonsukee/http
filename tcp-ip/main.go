package main

import "tcpip/tcpip"

func main() {
	// 宛先IPアドレスのMacアドレスを取得
	// TODO: 宛先IPアドレスを変更して試すとより良い！
	if err := tcpip.ArpRequest("en0", "192.168.1.1"); err != nil {
		panic(err)
	}
}
