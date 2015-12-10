/** Lucky 13 Proxy **/
package main

import (
	"github.com/jakemask/lucky13/defaults"
	"github.com/jakemask/lucky13/proxy"
	"log"
	"net"
	"net/rpc"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	// get proxy server up and running
	go serve(&wg)

	// connect to client
	client, err := rpc.DialHTTP("tcp", defaults.CLIENT_HOST+":"+defaults.CLIENT_PORT)
	if err != nil {
		log.Fatal("dialing error:", err)
	}

	var reply int
	if err := client.Call("Client.Send", []byte("lo"), &reply); err != nil {
		log.Fatal("client error:", err)
	}

	log.Println("waiting for connections")
	wg.Wait()
}

func serve(wg *sync.WaitGroup) {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", ":"+defaults.PROXY_PORT)
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	defer l.Close()

	log.Println("Listening on :" + defaults.PROXY_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		wg.Add(1)

		// Handle connections in a new goroutine.
		go handleRequest(wg, conn)
	}

}

// Handles incoming requests.
func handleRequest(wg *sync.WaitGroup, client net.Conn) {
	server, err := net.Dial("tcp", defaults.SERVER_HOST+":"+defaults.SERVER_PORT)
	if err != nil {
		log.Fatal("couldn't connect to server:", err)
	}

	proxy.Run(client, server, proxy.VerboseMITM, proxy.VerboseMITM)

	wg.Done()
}
