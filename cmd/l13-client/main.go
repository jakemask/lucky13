/** Lucky 13 Client **/
package main

import (
	"crypto/tls"
	"encoding/hex"
	"github.com/jakemask/lucky13/defaults"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Client int

func (t *Client) Send(msg []byte, reply *int) error {
	log.Println(hex.Dump(msg))

	config := tls.Config{
		InsecureSkipVerify: true,
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		},
	}
	conn, err := tls.Dial("tcp", defaults.PROXY_HOST+":"+defaults.PROXY_PORT, &config)
	if err != nil {
		log.Fatal("proxy connect error:", err)
	}
	defer conn.Close()

	if _, err := conn.Write(msg); err != nil {
		log.Printf("message error: ", err)
	}

	return nil
}

func main() {
	client := new(Client)

	rpc.Register(client)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+defaults.CLIENT_PORT)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	http.Serve(l, nil)
}
