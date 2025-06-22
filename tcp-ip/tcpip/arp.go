package tcpip

import (
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// 新しいARPリクエストパケットを作成
func NewArpRequest(srcIP net.IP, srcMAC net.HardwareAddr, targetIP net.IP) layers.ARP {
	return layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(srcMAC),
		SourceProtAddress: []byte(srcIP.To4()),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    []byte(targetIP.To4()),
	}
}

// 指定されたインターフェース名と宛先IPアドレスに対してARPリクエストを送信し，レスポンスを受信する関数
func Send(ifaceName string, targetIP net.IP) (*layers.ARP, error) {
	// パケットを送るために使用するインターフェース情報(例:有線LAN, 無線LAN)を取得
	// これは送信元が持つ情報だから送信元IPアドレスが登録されている
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, fmt.Errorf("インターフェースの取得に失敗: %v", err)
	}

	// インターフェースに割り当てられているIPアドレス一覧を取得
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("IPアドレスの取得に失敗: %v", err)
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
		return nil, fmt.Errorf("有効なIPv4アドレスが見つかりません")
	}

	// 送信元MACアドレスと宛先MACアドレス(まだわからないため全員へブロードキャスト)を設定
	srcMAC := iface.HardwareAddr
	dstMAC := net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	// パケットキャプチャ用のハンドルを開く(パケットを送受信するための窓口)
	handle, err := pcap.OpenLive(ifaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("pcapハンドルのオープンに失敗: %v", err)
	}
	defer handle.Close()

	// 実際に送信するパケットの中身を作成
	// ARP型のEthernetヘッダを作成
	eth := NewEthernet(srcMAC, dstMAC, "ARP")

	// ARPリクエストパケットを作成
	arp := NewArpRequest(srcIP, srcMAC, targetIP)

	// パケットをシリアライズ(バイト列に変換)
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	if err := gopacket.SerializeLayers(buf, opts, &eth, &arp); err != nil {
		return nil, fmt.Errorf("パケットのシリアライズに失敗: %v", err)
	}

	// ARPリクエストパケットを送信
	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		return nil, fmt.Errorf("パケットの送信に失敗: %v", err)
	}

	fmt.Printf("ARPリクエストを[%v]へ送信\n", targetIP)

	// ARPリプライを待ち受ける(3秒間)
	src := gopacket.NewPacketSource(handle, handle.LinkType())
	timeout := time.After(3 * time.Second)

	for {
		select {
		case packet := <-src.Packets():
			// ARPレイヤを抽出
			if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
				arpResponse := arpLayer.(*layers.ARP)
				// 宛先IPからのARPリプライか判定
				if net.IP(arpResponse.SourceProtAddress).Equal(targetIP) && arpResponse.Operation == layers.ARPReply {
					return arpResponse, nil
				}
			}
		case <-timeout:
			return nil, fmt.Errorf("タイムアウト: ARPリプライを受信できませんでした")
		}
	}
}

// ARPリクエストを実行する関数
func ArpRequest(ifaceName string, targetIPStr string) error {
	// 宛先IPアドレスを取得
	targetIP := net.ParseIP(targetIPStr)
	if targetIP == nil {
		return fmt.Errorf("無効なIPアドレス: %s", targetIPStr)
	}

	// ARPリクエストを送信
	arpReply, err := Send(ifaceName, targetIP)
	if err != nil {
		return fmt.Errorf("ARPリクエストの送信に失敗: %v", err)
	}

	// 結果を表示
	fmt.Printf(
		"ARPの結果: IP [%v]: MAC [%v]\n",
		net.IP(arpReply.SourceProtAddress),
		net.HardwareAddr(arpReply.SourceHwAddress),
	)

	return nil
}
