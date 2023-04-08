package gop0f

import (
	"fmt"
	"net"
)

type Client interface {
	Query(ip string) (*hostInfo, error)
}

type p0f struct {
	addr string
}

func New(socketAddress string) (*p0f, error) {
	return &p0f{
		addr: socketAddress,
	}, nil
}

func (p0f *p0f) Query(ip string) (*hostInfo, error) {
	// prepare data
	p := preparePacket(ip)

	conn, err := p0f.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// send data
	packet, err := send(conn, p)
	if err != nil {
		return nil, err
	}

	// parse received data
	result := &hostInfo{}
	packet.parse(result)

	return result, nil
}

func (p0f *p0f) connect() (net.Conn, error) {
	// todo implement connection pool
	conn, err := net.Dial("unix", p0f.addr)
	if err != nil {
		return nil, fmt.Errorf("initializing p0f client: %w", err)
	}
	return conn, nil
}
