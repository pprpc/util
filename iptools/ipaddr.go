package iptools

import (
	"encoding/binary"
	"errors"
	"net"
)

// IPv42long .
func IPv42long(ipv4 string) (uint32, error) {
	ip := net.ParseIP(ipv4)
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip), nil
}

// Long2ipv4 .
func Long2ipv4(ipv4i uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipv4i)
	ip := net.IP(ipByte)
	return ip.String()
}
