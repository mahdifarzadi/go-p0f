package gop0f

import (
	"fmt"
	"net"
)

type p0f struct {
	conn net.Conn
}

func New(socketAddress string) (*p0f, error) {
	conn, err := net.Dial("unix", socketAddress)
	if err != nil {
		return nil, fmt.Errorf("initializing p0f client: %w", err)
	}
	return &p0f{
		conn: conn,
	}, nil
}

func (p0f *p0f) Query(ip string) {
	// prepare data
	p := preparePacket(ip)
	fmt.Println(p)

	// send data
	packet, err := send(p0f.conn, p)
	fmt.Println(packet, err)

	// parse response to model

}
