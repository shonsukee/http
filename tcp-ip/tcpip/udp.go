package tcpip

import (
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// UDPパケットの構造体
type UDPPacket struct {
	Ethernet *layers.Ethernet
	IP       *layers.IPv4
	UDP      *layers.UDP
	Payload  []byte
}

// 新しいEthernetヘッダを作成
func NewUDPHeader(srcPort, dstPort uint16) *layers.UDP {
	return &layers.UDP{
		SrcPort: layers.UDPPort(srcPort),
		DstPort: layers.UDPPort(dstPort),
	}
}

// 新しいUDPパケットを作成
func NewUDPPacket(srcMAC, dstMAC net.HardwareAddr, srcIP, dstIP net.IP, srcPort, dstPort uint16, payload []byte) *UDPPacket {
	// Ethernetヘッダを作成
	ethernet := NewEthernet(srcMAC, dstMAC, "IPv4")

	// IPヘッダを作成
	ip := NewIPHeader(srcIP, dstIP)

	// UDPヘッダを作成
	udp := NewUDPHeader(srcPort, dstPort)

	// ペイロードを設定
	udp.SetNetworkLayerForChecksum(ip)

	return &UDPPacket{
		Ethernet: &ethernet,
		IP:       ip,
		UDP:      udp,
		Payload:  payload,
	}
}

// UDPパケットを送信する関数
func UdpSend(ifaceName string, dstIPStr string, srcPort, dstPort uint16, payload []byte) error {
	// インタフェース情報を取得
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return fmt.Errorf("インタフェースの取得に失敗: %v", err)
	}

	// インタフェースのIPアドレスを取得
	addrs, err := iface.Addrs()
	if err != nil {
		return fmt.Errorf("IPアドレスの取得に失敗: %v", err)
	}

	// IPv4アドレスを選択
	var srcIP net.IP
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				srcIP = ipnet.IP.To4()
				break
			}
		}
	}
	if srcIP == nil {
		return fmt.Errorf("有効なIPv4アドレスが見つかりません")
	}

	// 宛先IPアドレスを解析
	dstIP := net.ParseIP(dstIPStr)
	if dstIP == nil {
		return fmt.Errorf("無効な宛先IPアドレス: %s", dstIPStr)
	}

	// MACアドレスを取得(ARPを使用)
	srcMAC := iface.HardwareAddr

	// ARPを使用して宛先MACアドレスを取得
	arpReply, err := Send(ifaceName, dstIP)
	if err != nil {
		return fmt.Errorf("ARPリクエストに失敗: %v", err)
	}
	dstMAC := net.HardwareAddr(arpReply.SourceHwAddress)

	// UDPパケットを作成
	udpPacket := NewUDPPacket(srcMAC, dstMAC, srcIP, dstIP, srcPort, dstPort, payload)

	// パケットキャプチャ用のハンドルを開く
	handle, err := pcap.OpenLive(ifaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("pcapハンドルのオープンに失敗: %v", err)
	}
	defer handle.Close()

	// パケットをシリアライズ
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	if err := gopacket.SerializeLayers(buf, opts,
		udpPacket.Ethernet,
		udpPacket.IP,
		udpPacket.UDP,
		gopacket.Payload(udpPacket.Payload)); err != nil {
		return fmt.Errorf("パケットのシリアライズに失敗: %v", err)
	}

	// UDPパケットを送信
	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		return fmt.Errorf("パケットの送信に失敗: %v", err)
	}

	fmt.Printf("UDPパケットを[%v:%d]へ送信（送信元ポート: %d）\n", dstIP, dstPort, srcPort)
	return nil
}

// UDPリクエストを実行する関数
func UdpRequest(ifaceName string, dstIPStr string, srcPort, dstPort uint16, message string) error {
	payload := []byte(message)

	// UDPパケットを送信
	if err := UdpSend(ifaceName, dstIPStr, srcPort, dstPort, payload); err != nil {
		return fmt.Errorf("UDPパケットの送信に失敗: %v", err)
	}

	fmt.Printf("UDPメッセージ送信完了: %s\n", message)
	return nil
}
