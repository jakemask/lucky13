package proxy

import (
	"encoding/hex"
	"github.com/jakemask/lucky13/tlsparse"
	"log"
	"time"
)

type MITM func(time.Time, tlsparse.TLSHeader, []byte, string) (tlsparse.TLSHeader, []byte)

func NilMITM(_ time.Time, h tlsparse.TLSHeader, x []byte, _ string) (tlsparse.TLSHeader, []byte) {
	return h, x
}

var _ MITM = NilMITM

func VerboseMITM(t time.Time, hdr tlsparse.TLSHeader, msg []byte, desc string) (tlsparse.TLSHeader, []byte) {

	log.Printf("Arrived at %v", t)
	log.Printf("Header(%s): %x", desc, hdr)
	log.Println("Message:\n" + hex.Dump(msg))

	return hdr, msg
}

var _ MITM = VerboseMITM

// TODO this is super janky; there must be a better way
var start time.Time
var out bool

var Times []time.Duration

func ClientMITM(t time.Time, hdr tlsparse.TLSHeader, msg []byte, desc string) (tlsparse.TLSHeader, []byte) {

	if hdr.ContentType == 0x17 {
		// Application Message
		//log.Printf("Header(%s): %x", desc, hdr)
		//log.Println("Message:\n" + hex.Dump(msg))

		//msg = msg[0 : 288+16]
		//hdr.Length = 288 + 16

		out = true
		start = time.Now()
	}

	return hdr, msg
}

var _ MITM = ClientMITM

func ServerMITM(t time.Time, hdr tlsparse.TLSHeader, msg []byte, desc string) (tlsparse.TLSHeader, []byte) {
	if hdr.ContentType == 0x15 {
		if out {
			duration := t.Sub(start)
			Times = append(Times, duration)
			//log.Println("Duration: " + duration.String())
			out = false
		}
		// Alert Message
		//log.Printf("Header(%s): %x", desc, hdr)
		//log.Println("Message:\n" + hex.Dump(msg))

	}
	return hdr, msg
}
