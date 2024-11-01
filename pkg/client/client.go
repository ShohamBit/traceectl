package client

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// github.com/ShohamBit/traceectl holds the gRPC connection and service client.
const (
	// unix socket
	PROTOCOL_UNIX = "unix"
	PROTOCOL_TCP  = "tcp"
	SOCKET        = "/tmp/tracee.sock"
	DefaultIP     = "localhost"
	DefaultPort   = "4466"
)

type ServerInfo struct {
	ConnectionType string // Field to specify connection type (e.g., "unix" or "tcp")
	UnixSocketPath string // Path for the Unix socket, if using Unix connection
	ADDR           string // Address for the connection
}

// this function use grpc to connect the server
// it can connect to the server with tcp stream or unix socket
func connectToServer(serverInfo ServerInfo) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Use switch case to determine connection type
	var conn *grpc.ClientConn
	var err error

	switch serverInfo.ConnectionType {
	case PROTOCOL_UNIX:
		// Dial a Unix socket
		conn, err = grpc.NewClient(fmt.Sprintf("unix://%s", serverInfo.UnixSocketPath), opts...)
	case PROTOCOL_TCP:
		// Dial a TCP address
		address := fmt.Sprintf(serverInfo.ADDR)
		conn, err = grpc.NewClient(address, opts...)
	default:
		return nil, fmt.Errorf("unsupported connection type: %s", serverInfo.ConnectionType)
	}

	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
		return nil, err
	}
	return conn, nil
}
