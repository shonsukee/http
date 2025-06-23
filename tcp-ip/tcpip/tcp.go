package tcpip

import (
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// TCP接続の設定を表す構造体
type TCPIP struct {
	DestIP    string
	DestPort  uint16
	TcpFlag   string
	SeqNumber uint32
	AckNumber uint32
}

// TCP接続を管理する構造体
type TCPConnection struct {
	handle    *pcap.Handle
	ifaceName string
	srcIP     net.IP
	srcPort   uint16
	dstIP     net.IP
	dstPort   uint16
	seqNumber uint32
	ackNumber uint32
	srcMAC    net.HardwareAddr
	dstMAC    net.HardwareAddr
}

// 新しいTCP接続を作成
func NewTCP(ifaceName string, srcPort uint16) *TCPConnection {
	return &TCPConnection{
		ifaceName: ifaceName,
		srcPort:   srcPort,
		seqNumber: 1000, // 初期シーケンス番号
	}
}

// TCP接続を閉じる
func (t *TCPConnection) Close() error {
	if t.handle != nil {
		t.handle.Close()
	}
	return nil
}

// 新しいTCPヘッダを作成
func NewTcpHeader(srcPort, dstPort uint16, seq, ack uint32, flags string) *layers.TCP {
	tcp := &layers.TCP{
		SrcPort:    layers.TCPPort(srcPort),
		DstPort:    layers.TCPPort(dstPort),
		Seq:        seq,
		Ack:        ack,
		DataOffset: 5,
		Window:     65535,
		Checksum:   0,
		Urgent:     0,
	}

	// TCPフラグを設定
	switch flags {
	case "SYN":
		tcp.SYN = true
	case "ACK":
		tcp.ACK = true
	case "FIN":
		tcp.FIN = true
	case "SYNACK":
		tcp.SYN = true
		tcp.ACK = true
	case "FINACK":
		tcp.FIN = true
		tcp.ACK = true
	}

	return tcp
}

// 接続の初期設定を行う
func (t *TCPConnection) setupConnection(destIP string, destPort uint16) error {
	// インタフェース情報を取得
	iface, err := net.InterfaceByName(t.ifaceName)
	if err != nil {
		return fmt.Errorf("インタフェースの取得に失敗: %v", err)
	}

	// 送信元IPアドレスを取得
	addrs, err := iface.Addrs()
	if err != nil {
		return fmt.Errorf("IPアドレスの取得に失敗: %v", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				t.srcIP = ipnet.IP.To4()
				break
			}
		}
	}
	if t.srcIP == nil {
		return fmt.Errorf("有効なIPv4アドレスが見つかりません")
	}

	// 宛先IPアドレスを解析
	t.dstIP = net.ParseIP(destIP)
	if t.dstIP == nil {
		return fmt.Errorf("無効な宛先IPアドレス: %s", destIP)
	}

	t.dstPort = destPort
	t.srcMAC = iface.HardwareAddr

	// すでにMACアドレスがセットされていればARPをスキップ
	if t.dstMAC == nil {
		// ARPを使用して宛先MACアドレスを取得
		arpReply, err := Send(t.ifaceName, t.dstIP)
		if err != nil {
			return fmt.Errorf("ARPリクエストに失敗: %v", err)
		}
		t.dstMAC = net.HardwareAddr(arpReply.SourceHwAddress)
	}

	// パケットキャプチャ用のハンドルを開く
	handle, err := pcap.OpenLive(t.ifaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("pcapハンドルのオープンに失敗: %v", err)
	}
	t.handle = handle

	return nil
}

// TCPパケットを送信
func (t *TCPConnection) sendTCPPacket(tcp *layers.TCP, payload []byte) error {
	// Ethernetヘッダを作成
	ethernet := NewEthernet(t.srcMAC, t.dstMAC, "IPv4")

	// IPヘッダを作成
	ip := &layers.IPv4{
		SrcIP:    t.srcIP,
		DstIP:    t.dstIP,
		Protocol: layers.IPProtocolTCP,
		Version:  4,
		TTL:      64,
	}

	// TCPヘッダのチェックサムを設定
	tcp.SetNetworkLayerForChecksum(ip)

	// パケットをシリアライズ
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	layers := []gopacket.SerializableLayer{&ethernet, ip, tcp}
	if len(payload) > 0 {
		layers = append(layers, gopacket.Payload(payload))
	}

	if err := gopacket.SerializeLayers(buf, opts, layers...); err != nil {
		return fmt.Errorf("パケットのシリアライズに失敗: %v", err)
	}

	// TCPパケットを送信
	if err := t.handle.WritePacketData(buf.Bytes()); err != nil {
		return fmt.Errorf("パケットの送信に失敗: %v", err)
	}

	return nil
}

// TCPパケットを受信
func (t *TCPConnection) receiveTCPPacket(timeout time.Duration) (*layers.TCP, error) {
	src := gopacket.NewPacketSource(t.handle, t.handle.LinkType())
	timer := time.After(timeout)

	for {
		select {
		case packet := <-src.Packets():
			// TCPレイヤを抽出
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp := tcpLayer.(*layers.TCP)
				// 宛先ポートと送信元ポートを確認
				if tcp.DstPort == layers.TCPPort(t.srcPort) && tcp.SrcPort == layers.TCPPort(t.dstPort) {
					return tcp, nil
				}
			}
		case <-timer:
			return nil, fmt.Errorf("タイムアウト: TCPパケットを受信できませんでした")
		}
	}
}

// TCP接続を開始
func (t *TCPConnection) StartTCPConnection(tcpConfig TCPIP) (*TCPIP, error) {
	// 接続の初期設定
	if err := t.setupConnection(tcpConfig.DestIP, tcpConfig.DestPort); err != nil {
		return nil, err
	}

	// TCPヘッダを作成
	tcp := NewTcpHeader(t.srcPort, t.dstPort, t.seqNumber, t.ackNumber, tcpConfig.TcpFlag)

	// TCPパケットを送信
	if err := t.sendTCPPacket(tcp, nil); err != nil {
		return nil, err
	}

	fmt.Printf("TCP %sパケットを[%v:%d]へ送信\n", tcpConfig.TcpFlag, t.dstIP, t.dstPort)

	// レスポンスを受信(SYNの場合はSYN+ACKを待つ)
	if tcpConfig.TcpFlag == "SYN" {
		response, err := t.receiveTCPPacket(3 * time.Second)
		if err != nil {
			return nil, err
		}

		// シーケンス番号とACK番号を更新
		t.seqNumber = tcpConfig.SeqNumber + 1
		t.ackNumber = response.Seq + 1

		fmt.Printf("TCP SYN+ACKを受信: Seq=%d, Ack=%d\n", response.Seq, response.Ack)

		// ACKパケットを送信
		ackHeader := NewTcpHeader(t.srcPort, t.dstPort, t.seqNumber, t.ackNumber, "ACK")
		if err := t.sendTCPPacket(ackHeader, nil); err != nil {
			return nil, fmt.Errorf("ACKパケットの送信に失敗: %v", err)
		}
		fmt.Printf("TCP ACKパケットを[%v:%d]へ送信\n", t.dstIP, t.dstPort)

		return &TCPIP{
			DestIP:    tcpConfig.DestIP,
			DestPort:  tcpConfig.DestPort,
			TcpFlag:   "SYNACK",
			SeqNumber: t.seqNumber,
			AckNumber: t.ackNumber,
		}, nil
	}

	return &TCPIP{
		DestIP:    tcpConfig.DestIP,
		DestPort:  tcpConfig.DestPort,
		TcpFlag:   tcpConfig.TcpFlag,
		SeqNumber: t.seqNumber,
		AckNumber: t.ackNumber,
	}, nil
}
