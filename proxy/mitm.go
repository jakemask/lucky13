package proxy

import (
	"encoding/hex"
	"github.com/jakemask/lucky13/tlsparse"
	"log"
)

type MITM func(tlsparse.TLSHeader, []byte, string) (tlsparse.TLSHeader, []byte)

func NilMITM(h tlsparse.TLSHeader, x []byte, _ string) (tlsparse.TLSHeader, []byte) {
	return h, x
}

var _ MITM = NilMITM

func ClientMITM(hdr tlsparse.TLSHeader, msg []byte, desc string) (tlsparse.TLSHeader, []byte) {

	log.Printf("Header(%s): %x", desc, hdr)

	if hdr.ContentType == 0x17 {
		// Application Message
		log.Println("Message:\n" + hex.Dump(msg))
	}

	return hdr, msg
}

var _ MITM = ClientMITM

func VerboseMITM(hdr tlsparse.TLSHeader, msg []byte, desc string) (tlsparse.TLSHeader, []byte) {

	log.Printf("Header(%s): %x", desc, hdr)
	log.Println("Message:\n" + hex.Dump(msg))

	return hdr, msg
}

var _ MITM = VerboseMITM
