package gop0f

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	// address types
	ipv4Type = 0x04
	ipv6Type = 0x06

	// status
	statusBadQuery = 0x00
	statusOk       = 0x10
	statusNoMatch  = 0x20
)

var (
	//P0F_QUERY_MAGIC = [...]byte{0x1, 0x46, 0x30, 0x50} //0x50304601
	sendMagic = [...]byte{0x1, 0x46, 0x30, 0x50} //0x50304601
	recvMagic = [...]byte{0x2, 0x46, 0x30, 0x50} //0x50304602
)

type sendPacket struct {
	Magic    [4]byte  // Magic dword (0x50304601)
	AddrType byte     // Address type byte: 4 for IPv4, 6 for IPv6
	Data     [16]byte // 16 bytes of address data, network endian. IPv4 addresses should be aligned to the left
}

type recvPacket struct {
	Magic      [4]byte  // Magic dword (0x50304602)
	Status     [4]byte  // Status dword: 0x00 for 'bad query', 0x10 for 'OK', and 0x20 for 'no match'
	FirstSeen  uint32   // unix time (seconds) of first observation of the host
	LastSeen   uint32   // unix time (seconds) of most recent traffic
	TotalConn  uint32   // total number of connections seen
	UptimeMin  uint32   // calculated system uptime, in minutes. Zero if not known
	UpModDays  uint32   // uptime wrap-around interval, in days
	LastNat    [4]byte  // time of the most recent detection of IP sharing (NAT, load balancing, proxying). Zero if never detected
	LastChg    [4]byte  // time of the most recent individual OS mismatch (e.g., due to multiboot or IP reuse)
	Distance   [2]byte  // system distance (derived from TTL; -1 if no data)
	BadSw      byte     // p0f thinks the User-Agent or Server strings aren't accurate. The value of 1 means OS difference (possibly due to proxying), while 2 means an outright mismatch
	OsMatchQ   byte     // OS match quality: 0 for a normal match; 1 for fuzzy (e.g., TTL or DF difference); 2 for a generic signature; and 3 for both
	OsName     [32]byte // NUL-terminated name of the most recent positively matched OS. If OS not known, os_name[0] is NUL
	OsFlavor   [32]byte // OS version. May be empty if no data
	HttpName   [32]byte // most recent positively identified HTTP application (e.g. 'Firefox')
	HttpFlavor [32]byte // version of the HTTP application, if any
	LinkType   [32]byte // network link type, if recognized
	Language   [32]byte // system language, if recognized
}

func preparePacket(ipv4 string) *sendPacket {
	packet := &sendPacket{
		Magic:    sendMagic,
		AddrType: ipv4Type,
	}
	copy(packet.Data[:], ipv4toBytes(ipv4))
	return packet
}

func ipv4toBytes(ipv4 string) []byte {
	return net.ParseIP(ipv4).To4()
}

func send(conn net.Conn, packet *sendPacket) (*recvPacket, error) {
	// convert sendPacket to bytes array
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.BigEndian, packet); err != nil {
		return nil, err
	}
	data := buffer.Bytes()

	// check data length
	// todo

	// send data to p0f api
	_, err := conn.Write(data)
	if err != nil {
		return nil, err
	}

	// receive response data from p0f api
	recvData := make([]byte, 1048)
	recvLen, err := conn.Read(recvData[:])
	if err != nil {
		return nil, err
	}

	// convert recvData to recvPacket
	buf := bytes.NewReader(recvData[0:recvLen])
	fmt.Println("buf:", buf)
	recvPacket := &recvPacket{}
	err = binary.Read(buf, binary.BigEndian, recvPacket)
	if err != nil {
		return nil, err
	}

	// validate received packet
	// todo

	return recvPacket, nil
}
