package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	rrd "myapp/roundrobindialer"
)

// main demonstrates the usage of RoundRobinDialer by making HTTP requests
// to a target host with DNS round-robin resolution
func main() {
	target := os.Getenv("TARGET_DNS")
	if target == "" {
		target = "web.default.svc.cluster.local"
	}

	dialer := rrd.NewRoundRobinDialer(
		rrd.WithDialTimeout(3*time.Second),
		rrd.WithKeepAlive(10*time.Second),
		rrd.WithDNSTTL(30*time.Second),
	)

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext:       dialer.DialContext(),
			DisableKeepAlives: true,
		},
	}

	fmt.Printf("Starting requests to %s...\n", target)
	fmt.Println("----------------------------------------")

	for i := 0; i < 15; i++ {
		resp, err := client.Get("http://" + target)
		if err != nil {
			fmt.Printf("[%d] Request failed: %v\n", i+1, err)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Printf("[%d] Response from %s: %s", i+1, target, string(body))
		time.Sleep(500 * time.Millisecond)
	}
}
