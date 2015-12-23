package tlsparse

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

type Header struct {
	ContentType uint8
	Version     uint16
	Length      uint16
}

type Record struct {
	Header  Header
	Message []byte
}

func header(raw []byte) *Header {
	return &Header{
		raw[0],
		binary.BigEndian.Uint16(raw[1:3]),
		binary.BigEndian.Uint16(raw[3:5]),
	}
}

func (hdr *Header) Bytes() []byte {
	buf := make([]byte, 5)

	buf[0] = hdr.ContentType
	binary.BigEndian.PutUint16(buf[1:3], hdr.Version)
	binary.BigEndian.PutUint16(buf[3:5], hdr.Length)

	return buf
}

func (self *Record) Bytes() []byte {
	return append(self.Header.Bytes(), self.Message...)
}

func readHeader(conn io.Reader) (*Header, time.Time, error) {
	rawHdr := make([]byte, 5)
	hdrLen, err := conn.Read(rawHdr)
	t := time.Now()
	if err != nil {
		return nil, t, err
	}
	if hdrLen != 5 {
		return nil, t, errors.New("tlsparse: header not available")
	}

	return header(rawHdr), t, nil
}

func readMessage(conn io.Reader, hdr *Header) ([]byte, error) {
	rawMsg := make([]byte, hdr.Length)

	msgLen, err := conn.Read(rawMsg)

	// handle bad length error
	if msgLen != int(hdr.Length) {
		errMsg := "bad length"
		if err != nil {
			errMsg += ": " + err.Error()
		}
		return nil, errors.New(errMsg)
	} else if err == io.EOF {
		// suppress EOF if the message length was fine
		err = nil
	}

	return rawMsg, err
}

func ReadRecord(conn io.Reader) (*Record, time.Time, error) {
	hdr, t, err := readHeader(conn)
	if err != nil {
		return nil, t, err
	}

	msg, err := readMessage(conn, hdr)
	if err != nil {
		return nil, t, err
	}

	return &Record{*hdr, msg}, t, nil
}
