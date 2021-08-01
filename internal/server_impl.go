package internal

import (
	"log"
	"sync"

	pb "github.com/bartmika/tpoller-server/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type TelemetryServerImpl struct {
	mu            *sync.Mutex
	arduinoReader *ArduinoReader
	pb.TelemetryServer
}

func (s *TelemetryServerImpl) GetTimeSeriesData(in *empty.Empty, stream pb.Telemetry_GetTimeSeriesDataServer) error {
	// Since our device is a shared resource of the system, we need to coordinate
	// access accross different IPCs; therefore, we will protect our shared
	// resource with a `mutext` and grant access based on whether the `mutext`
	// was `locked` or `unlocked`.
	s.mu.Lock()
	defer s.mu.Unlock()

	data := s.arduinoReader.GetTimeSeriesData()
	for _, datum := range data {
		if err := stream.Send(datum); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
