package main

import (
	"crypto/tls"
	"encoding/hex"
	"github.com/jakemask/lucky13/defaults"
	"log"
	"net"
)

func main() {
	cer, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":"+defaults.SERVER_PORT, config)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("Recieved connection")

	tlscon, ok := conn.(*tls.Conn)
	if !ok {
		log.Fatal("couldn't cast to tls connection")
	}

	if err := tlscon.Handshake(); err != nil {
		log.Fatal("handshake error:", err)
	}
	log.Printf("%#v", tlscon.ConnectionState())

	buf := make([]byte, 1)
	var msg []byte

	for {
		n, err := conn.Read(buf)

		msg = append(msg, buf[:n]...)

		if err != nil {
			log.Println(err)
			break
		}

	}
	log.Println("Message:\n" + hex.Dump(msg))
}
