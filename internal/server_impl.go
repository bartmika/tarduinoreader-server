package internal

import (
	"log"

	pb "github.com/bartmika/tpoller-server/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type TelemetryServerImpl struct {
	arduinoReader *ArduinoReader
	pb.TelemetryServer
}

func (s *TelemetryServerImpl) PollTelemeter(in *empty.Empty, stream pb.Telemetry_PollTelemeterServer) error {
	data := s.arduinoReader.GetTimeSeriesData()
	for _, datum := range data {
		if err := stream.Send(datum); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
