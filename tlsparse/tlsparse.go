package tlsparse

import (
	"encoding/binary"
)

type TLSHeader struct {
	ContentType uint8
	Version     uint16
	Length      uint16
}

func Header(raw []byte) TLSHeader {
	return TLSHeader{
		raw[0],
		binary.BigEndian.Uint16(raw[1:3]),
		binary.BigEndian.Uint16(raw[3:5]),
	}
}

func (hdr TLSHeader) Bytes() []byte {
	buf := make([]byte, 5)

	buf[0] = hdr.ContentType
	binary.BigEndian.PutUint16(buf[1:3], hdr.Version)
	binary.BigEndian.PutUint16(buf[3:5], hdr.Length)

	return buf
}
