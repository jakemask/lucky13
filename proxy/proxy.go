package proxy

import (
	"github.com/jakemask/lucky13/tlsparse"
	"io"
	"log"
	"net"
	"sync"
)

const (
	DEBUG = false
)

func Run(client, server net.Conn, mitmClient, mitmServer MITM) {
	var wg sync.WaitGroup

	wg.Add(2)
	go pipe(&wg, client, server, mitmClient, "client -> server")
	go pipe(&wg, server, client, mitmServer, "server -> client")

	wg.Wait()

	log.Println("Closing proxy connections")

	client.Close()
	server.Close()
}

func pipe(wg *sync.WaitGroup, src, dst net.Conn, mitm MITM, desc string) {
	defer wg.Done()

	rawHdr := make([]byte, 5)
	rawMsg := make([]byte, 0xffff) // 64k buffer

	for {
		// read TLS record header
		hdrLen, err := src.Read(rawHdr)
		if hdrLen != 5 {
			if err != io.EOF {
				log.Println("bad length: ", err)
			}
			return
		}
		if DEBUG {
			log.Println(desc, "reading header")
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

		if DEBUG {
			log.Println(desc, "reading record")
		}

		// modify the TLS record
		newHdr, newMsg := mitm(tlsHdr, rawMsg[:msgLen], desc)

		// send the new TLS record
		if _, err = dst.Write(append(newHdr.Bytes(), newMsg...)); err != nil {
			log.Println("error writing: ", err)
			return
		}

		if DEBUG {
			log.Println(desc, "wrote record")
		}

		if read_err != nil {
			log.Println("read error: ", read_err)
			return
		}
	}
}
