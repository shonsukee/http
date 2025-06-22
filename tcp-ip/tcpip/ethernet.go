package tcpip

import (
	"net"

	"github.com/google/gopacket/layers"
)

// 新しいイーサネットフレームを作成
func NewEthernet(srcMAC net.HardwareAddr, dstMAC net.HardwareAddr, ethType string) layers.Ethernet {
	// 宛先，送信元のMacアドレスを設定
	ethernet := layers.Ethernet{
		DstMAC: dstMAC,
		SrcMAC: srcMAC,
	}

	// EtherTypeを設定
	switch ethType {
	case "IPv4":
		ethernet.EthernetType = layers.EthernetTypeIPv4
	case "ARP":
		ethernet.EthernetType = layers.EthernetTypeARP
	}
	return ethernet
}
