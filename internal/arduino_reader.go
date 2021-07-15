package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	pb "github.com/bartmika/tpoller-server/proto"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/tarm/serial"
)

const RX_BYTE = "1"
const SPARKFUN_WEATHER_SHIELD = "SPARKFUN-DEV-13956"

// The time-series data structure used to store all the data that will be
// returned by the `SparkFun Weather Shield` Arduino device.
//
// Please note that the value fields must be float32 as this was specified by
// the hardware manufacturar, if you use float64 you will get an error.
type SparkFunWeatherShieldTimeSeriesData struct {
	Status                 string  `json:"status,omitempty"`
	Runtime                int     `json:"runtime,omitempty"`
	Id                     int     `json:"id,omitempty"`
	HumidityValue          float32 `json:"humidity_value,omitempty"`
	HumidityUnit           string  `json:"humidity_unit,omitempty"`
	TemperatureValue       float32 `json:"temperature_primary_value,omitempty"`
	TemperatureUnit        string  `json:"temperature_primary_unit,omitempty"`
	PressureValue          float32 `json:"pressure_value,omitempty"`
	PressureUnit           string  `json:"pressure_unit,omitempty"`
	TemperatureBackupValue float32 `json:"temperature_secondary_value,omitempty"`
	TemperatureBackupUnit  string  `json:"temperature_secondary_unit,omitempty"`
	AltitudeValue          float32 `json:"altitude_value,omitempty"`
	AltitudeUnit           string  `json:"altitude_unit,omitempty"`
	IlluminanceValue       float32 `json:"illuminance_value,omitempty"`
	IlluminanceUnit        string  `json:"illuminance_unit,omitempty"`
	Timestamp              int64   `json:"timestamp,omitempty"`
}

// The abstraction of the `SparkFun Weather Shield` reader.
type ArduinoReader struct {
	serialPort     *serial.Port
	shieldHardware string
}

// Constructor used to intialize the serial reader designed to communicate
// with Arduino configured for the `SparkFun Weather Shield` settings.
func NewArduinoReader(devicePath string, arduinoShield string) *ArduinoReader {
	log.Printf("READER: Attempting to connect Arduino device...")
	c := &serial.Config{Name: devicePath, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	// DEVELOPERS NOTE:
	// The following code will warm up the Arduino device before we are
	// able to make calls to the external sensors.
	log.Printf("READER: Waiting for Arduino external sensors to warm up")
	ar := &ArduinoReader{
		serialPort:     s,
		shieldHardware: arduinoShield,
	}
	ar.GetTimeSeriesData()
	time.Sleep(5 * time.Second)
	ar.GetTimeSeriesData()
	time.Sleep(5 * time.Second)
	return ar
}

func (ar *ArduinoReader) GetTimeSeriesData() []*pb.TelemetryDatum {
	if ar.shieldHardware == SPARKFUN_WEATHER_SHIELD {

		d := ar.getSparkFunWeatherShieldData()
		if d == nil {
			return nil
		}

		// Generate our labels.
		labels := []*pb.TelemetryLabel{}
		labels = append(labels, &pb.TelemetryLabel{Name: "unit", Value: d.HumidityUnit})
		htsd := &pb.TelemetryDatum{
			Metric:    "humidity",
			Value:     float64(d.HumidityValue),
			Labels:    labels,
			Timestamp: &tspb.Timestamp{Seconds: d.Timestamp, Nanos: 0},
		}

		labels = []*pb.TelemetryLabel{}
		labels = append(labels, &pb.TelemetryLabel{Name: "unit", Value: d.TemperatureUnit})
		ttsd := &pb.TelemetryDatum{
			Metric:    "temperature",
			Value:     float64(d.TemperatureValue),
			Labels:    labels,
			Timestamp: &tspb.Timestamp{Seconds: d.Timestamp, Nanos: 0},
		}

		labels = []*pb.TelemetryLabel{}
		labels = append(labels, &pb.TelemetryLabel{Name: "unit", Value: d.PressureUnit})
		ptsd := &pb.TelemetryDatum{
			Metric:    "pressure",
			Value:     float64(d.PressureValue),
			Labels:    labels,
			Timestamp: &tspb.Timestamp{Seconds: d.Timestamp, Nanos: 0},
		}

		labels = []*pb.TelemetryLabel{}
		labels = append(labels, &pb.TelemetryLabel{Name: "unit", Value: d.TemperatureBackupUnit})
		tbtsd := &pb.TelemetryDatum{
			Metric:    "temperature_backup",
			Value:     float64(d.TemperatureBackupValue),
			Labels:    labels,
			Timestamp: &tspb.Timestamp{Seconds: d.Timestamp, Nanos: 0},
		}

		labels = []*pb.TelemetryLabel{}
		labels = append(labels, &pb.TelemetryLabel{Name: "unit", Value: d.AltitudeUnit})
		atsd := &pb.TelemetryDatum{
			Metric:    "altitude",
			Value:     float64(d.AltitudeValue),
			Labels:    labels,
			Timestamp: &tspb.Timestamp{Seconds: d.Timestamp, Nanos: 0},
		}

		labels = []*pb.TelemetryLabel{}
		labels = append(labels, &pb.TelemetryLabel{Name: "unit", Value: d.IlluminanceUnit})
		itsd := &pb.TelemetryDatum{
			Metric:    "illuminance",
			Value:     float64(d.IlluminanceValue),
			Labels:    labels,
			Timestamp: &tspb.Timestamp{Seconds: d.Timestamp, Nanos: 0},
		}

		results := make([]*pb.TelemetryDatum, 6)
		results[0] = htsd
		results[1] = ttsd
		results[2] = ptsd
		results[3] = tbtsd
		results[4] = atsd
		results[5] = itsd
		return results
	}
	return nil
}

// Function returns the JSON data of the instrument readings from our Arduino
// device configured for the `SparkFun Weather Shield` settings.
func (ar *ArduinoReader) getSparkFunWeatherShieldData() *SparkFunWeatherShieldTimeSeriesData {
	// DEVELOPERS NOTE:
	// (1) The external device (Arduino) is setup to standby idle until it
	//     receives a poll request from this code, once a poll request has
	//     been submitted then all the sensors get polled and their data is
	//     returned.
	// (2) Please look at the following code to understand how the external
	//     device works in: <TODO ADD>
	//
	// (3) The reason for design is as follows:
	//     (a) The external device does not have a real-time clock
	//     (b) We don't want to add any real-time clock shields because
	//         extra hardware means it costs more.
	//     (c) We don't want to write complicated code of synching time
	//         from this code because it will make the code complicated.
	//     (d) Therefore we chose to make sensor polling be event based
	//         and this code needs to send a "poll request".

	// STEP 1:
	// We need to send a single byte to the external device (Arduino) which
	// will trigger a polling event on all the sensors.
	n, err := ar.serialPort.Write([]byte(RX_BYTE))
	if err != nil {
		log.Fatal(err)
	}

	// STEP 2:
	// The external device will poll the device, we need to make our main
	// runtime loop to be blocked so we wait until the device finishes and
	// returns all the sensor measurements.
	buf := make([]byte, 1028)
	n, err = ar.serialPort.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	// STEP 3:
	// Check to see if ANY data was returned from the external device, if
	// there was then we load up the string into a JSON object.
	var tsd SparkFunWeatherShieldTimeSeriesData
	err = json.Unmarshal(buf[:n], &tsd)
	if err != nil {
		return nil
	}
	tsd.Timestamp = time.Now().Unix()
	return &tsd
}

// Function used to print to the console the time series data.
func PrettyPrintTimeSeriesData(tsd *SparkFunWeatherShieldTimeSeriesData) {
	fmt.Println("Status: ", tsd.Status)
	fmt.Println("Runtime: ", tsd.Runtime)
	fmt.Println("Status: ", tsd.Id)
	fmt.Println("HumidityValue: ", tsd.HumidityValue)
	fmt.Println("HumidityUnit: ", tsd.HumidityUnit)
	fmt.Println("TemperatureValue: ", tsd.TemperatureValue)
	fmt.Println("TemperatureUnit: ", tsd.TemperatureUnit)
	fmt.Println("PressureValue: ", tsd.PressureValue)
	fmt.Println("PressureUnit: ", tsd.PressureUnit)
	fmt.Println("TemperatureBackupValue: ", tsd.TemperatureBackupValue)
	fmt.Println("TemperatureBackupUnit: ", tsd.TemperatureBackupUnit)
	fmt.Println("AltitudeValue: ", tsd.AltitudeValue)
	fmt.Println("AltitudeUnit: ", tsd.AltitudeUnit)
	fmt.Println("IlluminanceValue: ", tsd.IlluminanceValue)
	fmt.Println("IlluminanceUnit: ", tsd.IlluminanceUnit)
	fmt.Println("Timestamp: ", tsd.Timestamp)
}
