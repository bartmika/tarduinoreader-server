package internal

import (
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/bartmika/tpoller-server/proto"
)

type TArduinoReaderServerImpl struct {
	arduinoReader *ArduinoReader
	pb.TPollerServer
}

func (s *TArduinoReaderServerImpl) PollTimeSeriesData(in *empty.Empty, stream pb.TPoller_PollTimeSeriesDataServer) error {
	data := s.arduinoReader.GetTimeSeriesData()
	for _, datum := range data {
		if err := stream.Send(datum); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
