package proxy

import (
	"github.com/jakemask/lucky13/tlsparse"
	"log"
	"net"
	"sync"
	"time"
)

const (
	VERBOSE_HANDSHAKE = false
)

const (
	recordTypeHandshake uint8 = 22
	recordTypeCCS       uint8 = 20
)

type ConnPair struct {
	client, server net.Conn
}

func (self *ConnPair) Handshake() {
	var wg sync.WaitGroup

	wg.Add(2)
	go pipe(&wg, self.client, self.server, "client -> server")
	go pipe(&wg, self.server, self.client, "server -> client")

	wg.Wait()

}

func pipe(wg *sync.WaitGroup, src, dst net.Conn, desc string) {
	defer wg.Done()

	seenCCS := false

	for {
		// default to nil mitm, verbose if debug
		mitm := NilMITM
		if VERBOSE_HANDSHAKE {
			mitm = VerboseMITM
		}

		record, _, err := send(src, dst, mitm)
		if err != nil {
			log.Printf("(%s) couldn't forward: %v", desc, err)
		}

		// check if we've seen the ChangeCipherSpec record
		seenCCS = seenCCS || record.Header.ContentType == recordTypeCCS

		// exit loop if we've seen the CCS and one more record (Finished)
		bail := seenCCS && record.Header.ContentType == recordTypeHandshake

		if bail {
			return
		}
	}
}

func send(src, dst net.Conn, mitm MITM) (*tlsparse.Record, time.Time, error) {
	record, t, err := tlsparse.ReadRecord(src)
	if err != nil {
		log.Println("error reading: ", err)
		return record, t, err
	}

	newRecord := mitm(record)

	_, err = dst.Write(newRecord.Bytes())
	if err != nil {
		log.Println("error writing: ", err)
		return newRecord, t, err
	}

	return newRecord, t, nil
}

func (self *ConnPair) Forward(mitm MITM) time.Time {
	send(self.client, self.server, mitm)
	return time.Now()
}

func (self *ConnPair) Receive(mitm MITM) time.Time {
	_, t, _ := send(self.server, self.client, mitm)
	return t
}

func (self *ConnPair) Close() {
	self.client.Close()
	self.server.Close()
}
