package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"

	pb "github.com/bartmika/tpoller-server/proto"
)

func init() {
	// The following are optional and will have defaults placed when missing.
	getDataCmd.Flags().IntVarP(&port, "port", "p", 50051, "The port of our server.")
	rootCmd.AddCommand(getDataCmd)
}

func doGetData() {
	// Set up a direct connection to the gRPC server.
	conn, err := grpc.Dial(
		fmt.Sprintf(":%v", port),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// Set up our protocol buffer interface.
	client := pb.NewTelemetryClient(conn)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Perform our gRPC request.
	stream, err := client.GetTimeSeriesData(ctx, &empty.Empty{})
	if err != nil {
		log.Fatalf("could not poll time series data: %v", err)
	}

	// Handle our stream of data from the server.
	for {
		dataPoint, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error with stream: %v", err)
		}

		// Print out the gRPC response.
		log.Printf("Server Response: %s", dataPoint)
	}
}

var getDataCmd = &cobra.Command{
	Use:   "poll",
	Short: "Poll data from the gRPC server",
	Long:  `Connect to the gRPC server and poll the time series data. Command used to test out that the server is running.`,
	Run: func(cmd *cobra.Command, args []string) {
		doGetData()
	},
}
