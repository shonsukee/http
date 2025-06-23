package tcpip

import (
	"net"

	"github.com/google/gopacket/layers"
)

func NewIPHeader(srcIP, dstIP net.IP) *layers.IPv4 {
	return &layers.IPv4{
		Version:  4,
		IHL:      5,
		Protocol: layers.IPProtocolUDP,
		SrcIP:    srcIP,
		DstIP:    dstIP,
	}
}
