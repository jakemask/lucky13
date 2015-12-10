/** Lucky 13 Proxy **/
package main

import (
	"github.com/jakemask/lucky13/defaults"
	"github.com/jakemask/lucky13/proxy"
	"log"
	"net"
	"net/rpc"
	"sort"
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

	//m2 := true //false should be slower

	msg := []byte("12345678901234567890")

	test(client, &wg, msg)
	msg = append(msg, msg...)
	test(client, &wg, msg)
}

func test(client *rpc.Client, wg *sync.WaitGroup, msg []byte) {
	for i := 0; i < 128; i++ {
		var reply int

		wg.Add(1)
		if err := client.Call("Client.Send", msg, &reply); err != nil {
			log.Fatal("client error:", err)
		}
	}

	log.Println("waiting for connections", wg)
	wg.Wait()

	sum := int64(0)
	var ns []int
	for _, v := range proxy.Times {
		sum += v.Nanoseconds()
		ns = append(ns, int(v.Nanoseconds()))
	}
	avg := sum / int64(len(proxy.Times))
	sort.Ints(ns)
	median := ns[len(ns)/2]

	log.Println("Took an average of ", float64(avg)/1000, " microseconds")
	log.Println("Took a median of ", float64(median)/1000, " microseconds")

	proxy.Times = nil

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

		// Handle connections in a new goroutine.
		handleRequest(conn)

		wg.Done()
	}

}

// Handles incoming requests.
func handleRequest(client net.Conn) {
	server, err := net.Dial("tcp", defaults.SERVER_HOST+":"+defaults.SERVER_PORT)
	if err != nil {
		log.Fatal("couldn't connect to server:", err)
	}

	proxy.Run(client, server, proxy.ClientMITM, proxy.ServerMITM)
}
