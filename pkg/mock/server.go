package mock

import (
	"fmt"
	"net"
	"os"

	"github.com/ShohamBit/traceectl/pkg/client"
	pb "github.com/aquasecurity/tracee/api/v1beta1"

	"google.golang.org/grpc"
)

var (
	ExpectedVersion string            = "v0.22.0-15-gd09d7fca0d" // Match the output format
	serverInfo      client.ServerInfo = client.ServerInfo{
		ADDR:           client.DefaultIP + ":" + client.DefaultPort,
		UnixSocketPath: client.SOCKET,
	}
)

// MockServiceServer implements the gRPC server interface for testing
type MockServiceServer struct {
	pb.UnimplementedTraceeServiceServer // Embed the unimplemented server
}

// CreateMockServer initializes the gRPC server and binds it to a Unix socket listener
func CreateMockServer() (*grpc.Server, net.Listener, error) {
	//check for unix socket
	if _, err := os.Stat(serverInfo.UnixSocketPath); err == nil {
		err := os.Remove(serverInfo.UnixSocketPath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to cleanup gRPC listening address (%s): %v", serverInfo.UnixSocketPath, err)
		}
	}

	// Create the Unix socket listener
	listener, err := net.Listen("unix", serverInfo.UnixSocketPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Unix socket listener: %v", err)
	}

	// Create a new gRPC server
	server := grpc.NewServer()

	return server, listener, nil
}

func StartMockServiceServer() (*grpc.Server, error) {
	mockServer, listener, err := CreateMockServer()
	if err != nil {
		return nil, fmt.Errorf("failed to create mock server: %v", err)
	}
	pb.RegisterTraceeServiceServer(mockServer, &MockServiceServer{})

	// Start serving in a goroutine
	go func() {
		if err := mockServer.Serve(listener); err != nil {
			fmt.Printf("gRPC server failed: %v\n", err)
		}
	}()
	return mockServer, nil
}