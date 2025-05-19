# Round Robin Dialer

A Go library that implements DNS round-robin connection distribution for HTTP clients. This library helps distribute network connections across multiple IP addresses when a hostname resolves to multiple IPs, providing better load balancing and high availability.

## Features

- DNS round-robin resolution
- Configurable DNS TTL
- Customizable connection timeouts
- Keep-alive connection support
- Thread-safe implementation

## Installation

```bash
go get github.com/yourusername/roundrobin-dialer
```

## Usage

### Basic Usage

```go
import (
    "net/http"
    "time"
    rrd "github.com/yourusername/roundrobin-dialer"
)

// Create a new dialer with default settings
dialer := rrd.NewRoundRobinDialer()

// Create an HTTP client with the round-robin dialer
client := &http.Client{
    Transport: &http.Transport{
        DialContext: dialer.DialContext(),
    },
}
```

### Custom Configuration

```go
dialer := rrd.NewRoundRobinDialer(
    rrd.WithDialTimeout(3*time.Second),
    rrd.WithKeepAlive(10*time.Second),
    rrd.WithDNSTTL(30*time.Second),
)
```

## Configuration Options

- `WithDNSTTL`: Sets the DNS cache TTL duration (default: 30s)
- `WithKeepAlive`: Sets the keep-alive duration for connections (default: 10s)
- `WithDialTimeout`: Sets the connection timeout duration (default: 3s)

## Example

The project includes a simple example that demonstrates the usage of the RoundRobinDialer. To run the example:

```bash
# Set the target hostname (optional)
export TARGET_DNS=your-target-hostname

# Run the example
go run main.go
```

The example will make 15 HTTP requests to the target hostname, demonstrating how connections are distributed across different IP addresses.

## Docker Support

The project includes a Dockerfile for containerization. To build and run the container:

```bash
# Build the container
docker build -t roundrobin-dialer .

# Run the container
docker run -e TARGET_DNS=your-target-hostname roundrobin-dialer
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 