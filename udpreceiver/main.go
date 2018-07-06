package main

import (
	"flag"
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

var (
	listenAddr = flag.String("listen-addr", ":10000", "listen addr")
	batchSize  = flag.Int("batch-size", 1000, "batch size")
)

func main() {
	flag.Parse()
	ra, err := net.ResolveUDPAddr("udp", *listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", ra)
	if err != nil {
		log.Fatal(err)
	}
	pconn := ipv4.NewPacketConn(conn)

	rb := make([]ipv4.Message, *batchSize)
	for i := 0; i < *batchSize; i++ {
		rb[i].Buffers = [][]byte{make([]byte, 1500)}
	}
	count := 0
	bytes := 0
	for {
		n, err := pconn.ReadBatch(rb, 0)
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range rb[:n] {
			// m.Buffers[0][:m.N]
			bytes += m.N
		}
		count += n
	}
}
