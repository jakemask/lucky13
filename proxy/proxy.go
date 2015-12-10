package proxy

import (
	"github.com/jakemask/lucky13/tlsparse"
	"log"
	"net"
	"sync"
	"time"
)

func Run(client, server net.Conn, mitmClient, mitmServer MITM) {
	//defer log.Println("Closing proxy connections")
	defer client.Close()
	defer server.Close()

	var wg sync.WaitGroup

	wg.Add(2)
	go pipe(&wg, client, server, mitmClient, "client -> server")
	go pipe(&wg, server, client, mitmServer, "server -> client")

	wg.Wait()

}

func pipe(wg *sync.WaitGroup, src, dst net.Conn, mitm MITM, desc string) {
	defer wg.Done()

	rawHdr := make([]byte, 5)
	rawMsg := make([]byte, 0xffff) // 64k buffer

	for {
		// read TLS record header
		hdrLen, err := src.Read(rawHdr)
		arrival := time.Now()
		if hdrLen != 5 {
			if hdrLen != 0 {
				log.Println("bad length: ", hdrLen, err)
			}
			return
		}

		tlsHdr := tlsparse.Header(rawHdr)

		if err != nil {
			log.Println("error reading:", err)
			return
		}

		// read TLS record
		msgLen, read_err := src.Read(rawMsg[:tlsHdr.Length])
		if msgLen != int(tlsHdr.Length) {
			log.Println("bad length: ", err)
			return
		}

		// modify the TLS record
		newHdr, newMsg := mitm(arrival, tlsHdr, rawMsg[:msgLen], desc)

		// send the new TLS record
		_, err = dst.Write(append(newHdr.Bytes(), newMsg...))
		if err != nil {
			//log.Println("error writing: ", err)
			return
		}

		if read_err != nil {
			log.Println("read error: ", read_err)
			return
		}
	}
}
