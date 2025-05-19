package roundrobindialer

import (
	"context"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

// RoundRobinDialer implements a DNS round-robin dialer that distributes connections
// across multiple IP addresses for a given hostname
type RoundRobinDialer struct {
	mu          sync.Mutex
	next        int
	dnsTTL      time.Duration
	dialTimeout time.Duration
	keepAlive   time.Duration
}

// NewRoundRobinDialer creates a new RoundRobinDialer with optional configuration
func NewRoundRobinDialer(options ...func(*RoundRobinDialer)) *RoundRobinDialer {
	d := &RoundRobinDialer{
		dnsTTL:      30 * time.Second,
		dialTimeout: 3 * time.Second,
		keepAlive:   10 * time.Second,
	}
	for _, opt := range options {
		opt(d)
	}
	return d
}

// DialContext returns a dialer function that implements round-robin DNS resolution
func (d *RoundRobinDialer) DialContext() func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			port = "80"
			host = addr
		}

		ips, err := net.LookupIP(host)
		if err != nil || len(ips) == 0 {
			return nil, fmt.Errorf("DNS lookup failed for %s: %v", host, err)
		}

		sort.Slice(ips, func(i, j int) bool {
			return ips[i].String() < ips[j].String()
		})

		fmt.Printf("\nResolved IPs for %s: %v\n", host, ips)

		d.mu.Lock()
		ip := ips[d.next%len(ips)]
		d.next++
		d.mu.Unlock()

		fmt.Printf("Dialing IP: %s\n", ip.String())

		address := net.JoinHostPort(ip.String(), port)
		dialer := &net.Dialer{
			Timeout:   d.dialTimeout,
			KeepAlive: d.keepAlive,
		}
		return dialer.DialContext(ctx, network, address)
	}
}

// WithDNSTTL sets the DNS cache TTL duration
func WithDNSTTL(dur time.Duration) func(*RoundRobinDialer) {
	return func(s *RoundRobinDialer) {
		s.dnsTTL = dur
	}
}

// WithKeepAlive sets the keep-alive duration for connections
func WithKeepAlive(dur time.Duration) func(*RoundRobinDialer) {
	return func(s *RoundRobinDialer) {
		s.keepAlive = dur
	}
}

// WithDialTimeout sets the connection timeout duration
func WithDialTimeout(dur time.Duration) func(*RoundRobinDialer) {
	return func(s *RoundRobinDialer) {
		s.dialTimeout = dur
	}
}
