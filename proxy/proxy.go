package proxy

import (
	"crypto/tls"
	"encoding/hex"
	"log"
	"net"
	"time"
)

const (
	DEBUG = false
)

type Proxy struct {
	pairs    chan ConnPair
	listener net.Listener
	config   Config
}

type Config struct {
	ProxyPort  string
	ServerHost string
	ServerPort string
}

func Serve(config Config) *Proxy {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", ":"+config.ProxyPort)
	if err != nil {
		log.Fatal("listen error: ", err)
	}

	proxy := &Proxy{
		pairs:    make(chan ConnPair),
		listener: l,
		config:   config,
	}

	log.Println("Listening on :" + config.ProxyPort)
	go proxy.listen()

	return proxy
}

func (self *Proxy) listen() {
	defer self.listener.Close()

	for {
		// Listen for an incoming connection.
		conn, err := self.listener.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		if DEBUG {
			log.Printf("recieved connection %v", conn)
		}

		// connect to remote server
		server, err := net.Dial("tcp", self.config.ServerHost+":"+self.config.ServerPort)
		if err != nil {
			log.Fatal("couldn't connect to server:", err)
		}

		if DEBUG {
			log.Printf("connected to server %v", server)
		}

		// proxy the handshake then hand off the connection
		pair := ConnPair{conn, server}
		pair.Handshake()

		if DEBUG {
			log.Printf("handshake complete")
		}

		self.pairs <- pair

		if DEBUG {
			log.Printf("pair retrieved, looping")
		}
	}
}

func (self *Proxy) Send(msg []byte, mitm MITM) time.Duration {
	if DEBUG {
		log.Println("Message:\n" + hex.Dump(msg))
	}

	config := tls.Config{
		InsecureSkipVerify: true,
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		},
	}
	conn, err := tls.Dial("tcp", "localhost:"+self.config.ProxyPort, &config)
	if err != nil {
		log.Fatal("proxy connect error:", err)
	}
	defer conn.Close()

	if DEBUG {
		log.Printf("made tls connection")
	}

	pair := <-self.pairs
	defer pair.Close()

	if DEBUG {
		log.Printf("retrieved pair %v", pair)
	}

	if _, err := conn.Write(msg); err != nil {
		log.Printf("message error: ", err)
	}

	if DEBUG {
		log.Printf("tls message sent")
	}

	start := pair.Forward(mitm)

	if DEBUG {
		log.Printf("message forwarded at %v", start)
	}

	end := pair.Receive(NilMITM) //TODO actually check that it's an error?

	if DEBUG {
		log.Printf("reply received at %v", end)
	}

	return end.Sub(start)
}
