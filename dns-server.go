package main

import (
	"fmt"
	"net"
)

// <domain-name> is a domain name represented as a series of labels, and
// terminated by a label with zero length.

func main() {
	address, _ := net.ResolveUDPAddr("udp", "localhost:53")

	in_conn, err := net.ListenUDP("udp", address)
	if err != nil {
		panic(fmt.Errorf("error %s starting udp in_conn on %s", err.Error(), address))
	}

	out_conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		panic(fmt.Errorf("error %s starting udp out_conn on %s", err.Error(), address))
	}

	dns_resolver := net.UDPAddr{
		IP:   net.IPv4(0x08, 0x08, 0x08, 0x08),
		Port: 53,
		Zone: "udp",
	}
	for {
		var query [512]byte

		_, _, _, addr, _ := in_conn.ReadMsgUDP(query[:], nil)
		out_conn.WriteMsgUDP(query[:], nil, &dns_resolver)

		var response [512]byte

		out_conn.ReadMsgUDP(response[:], nil)
		fmt.Println(retrieve_domain(response[12:]))
		in_conn.WriteMsgUDP(response[:], nil, addr)
	}
}

// The total length of a domain name (i.e., label octets and
// label length octets) is restricted to 255 octets or less.
func retrieve_domain(msg []byte) string {
	var domain []byte

	for {
		size := msg[0]
		if size == 0 {
			break
		}
		size++
		tmp := make([]byte, size)
		copy(tmp, msg[1:])
		tmp[len(tmp)-1] = 0x2e
		domain = append(domain, tmp...)
		msg = msg[size:]
	}

	return string(domain)
}
