package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

var (
	dstAddr    = flag.String("dst", "127.0.0.1:10000", "dst addr")
	count      = flag.Int("count", 10000, "packet count")
	batchSize  = flag.Int("batch-size", 1000, "batch size")
	packetSize = flag.Int("packet-size", 32, "packet size")
)

func main() {
	flag.Parse()
	ra, err := net.ResolveUDPAddr("udp", *dstAddr)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, ra)
	if err != nil {
		log.Fatal(err)
	}
	pconn := ipv4.NewPacketConn(conn)

	wb := make([]ipv4.Message, *batchSize)
	for i := 0; i < *batchSize; i++ {
		wb[i].Addr = ra
		wb[i].Buffers = [][]byte{make([]byte, *packetSize)}
	}
	c := 0
	start := time.Now()
	for {
		n, err := pconn.WriteBatch(wb, 0)
		if err != nil {
			log.Fatal(err)
		}
		c += n
		if c >= *count {
			break
		}
	}
	end := time.Now()
	sub := end.Sub(start)
	fmt.Printf("send %d packets in %s\n", c, sub)
	fmt.Printf("%f pps\n", float64(c)/sub.Seconds())
}
