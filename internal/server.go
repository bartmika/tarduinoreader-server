package internal

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/bartmika/tpoller-server/proto"
)

type TReaderServer struct {
	port              int
	arduinoDevicePath string
	arduinoShield     string
	arduinoReader     *ArduinoReader
	grpcServer        *grpc.Server
}

func New(arduinoDevicePath string, arduinoShield string, port int) *TReaderServer {
	return &TReaderServer{
		port:              port,
		arduinoDevicePath: arduinoDevicePath,
		arduinoShield:     arduinoShield,
		arduinoReader:     nil,
		grpcServer:        nil,
	}
}

// Function will consume the main runtime loop and run the business logic
// of the application.
func (s *TReaderServer) RunMainRuntimeLoop() {
	// Open a TCP server to the specified localhost and environment variable
	// specified port number.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Establish our device connection.
	arduinoReader := NewArduinoReader(s.arduinoDevicePath, s.arduinoShield)

	// Initialize our gRPC server using our TCP server.
	grpcServer := grpc.NewServer()

	// Save reference to our application state.
	s.grpcServer = grpcServer
	s.arduinoReader = arduinoReader

	// For debugging purposes only.
	log.Printf("gRPC server is running.")

	// Block the main runtime loop for accepting and processing gRPC requests.
	pb.RegisterTelemetryServer(grpcServer, &TelemetryServerImpl{
		// DEVELOPERS NOTE:
		// We want to attach to every gRPC call the following variables...
		arduinoReader: arduinoReader,
	})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Function will tell the application to stop the main runtime loop when
// the process has been finished.
func (s *TReaderServer) StopMainRuntimeLoop() {
	log.Printf("Starting graceful shutdown now...")

	s.arduinoReader = nil

	// Finish any RPC communication taking place at the moment before
	// shutting down the gRPC server.
	s.grpcServer.GracefulStop()
}
