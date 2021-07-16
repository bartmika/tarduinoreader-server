# Time-Series Data Reader and gRPC Server

Read **time-series data**, using `serial communication over USB`, for your externally connected *Arduino device* and attached *Arduino shield*; afterwords, share your data through [`gRPC`](https://github.com/bartmika/tpoller-server/blob/master/proto/telemetry.proto).

Currently this server only supports the following shield:

* [SparkFun Weather Shield (DEV-13956)](https://github.com/sparkfun/Weather_Shield).

Example:
![SparkFun Weather Shield](https://github.com/bartmika/sparkfunweathershield-arduino/blob/master/media/red_mulberries_germination_with_sparkfun_weather_shield.jpg?raw=true)
*SparkFun Weather Shield monitoring dutifully a red mulberry seedling, the code is powered by this code repository* - read more via [this blog post](https://bartlomiejmika.com/post/2021/red-mulberry-growlog-2/).

## Background
### The Problem

Imagine having purchased multiple measuring instruments for an IoT project, for example: pressure monitor, temperature probe, multi-sensor reader, etc. Each sensor provides different measured data - **Is it possible to create a simple *get data from device* interface to abstract all implantation detail and provide an easy-to-use API?**

Imagine having built an application which can provide an easy-to-use API for multiple hardware sensors. **How do you make the API accessible between *one* or *many* applications running concurrently? How do you make this API accessible either *locally* and or *remotely*?**

### The Solution

`treader-server` tries to provide an easy-to-use interface for you to connect many interprocess communication requests either locally on your computer or remotely across the internet or your network.

### What do I need?
To run this code, you need the following:

1. You need an Arduino device.
2. You need to have the Arduino device running either one of the following code:

  * [sparkfunweathershield-arduino](https://github.com/bartmika/sparkfunweathershield-arduino)
3. The computer running this server needs to be connected to the Arduino device over `serial communication` USB port.
4. The computer running this server needs to know the port address of your externally connected Arduino. Please find out what USB port your external device is connected on. Note: please replace ``/dev/cu.usbmodem14201`` with the value on your machine, a Raspberry Pi would most likely have the value ``/dev/ttyACM0``.

### How does this work?

This server starts by connecting to the Arduino device, warms up sensors if necessary and runs continuously waiting for you to make gRPC requests.

When you make a gRPC request, this server will pull data from the external device and send you back a gRPC response. You do what you want with the data.

### What is the gRPC service definition?

Your application must implement the `Telemetry` service definition found [here](https://github.com/bartmika/tpoller-server/blob/master/proto/telemetry.proto). The service definition code snippet is as follows:

```proto
service Telemetry {
    rpc PollTelemeter (google.protobuf.Empty) returns (stream TelemetryDatum) {}
}

message TelemetryLabel {
    string name = 1;
    string value = 2;
}

message TelemetryDatum {
    string metric = 1;
    repeated TelemetryLabel labels = 2;
    double value = 3;
    google.protobuf.Timestamp timestamp = 4;
}
```

### Why did you choose Arduino?
The Arduino platform has a wonderful ecosystem of open-source hardware with libraries. The goal is to take advantage of the libraries the hardware manufacturers wrote and not worry about the complicated implementation details; nor conflicting non-open source licensing agreements.

### Why should I use it?
* You want to focus on writing software or web-applications utilizing IoT sensors, not focus low-level hardware code.

* You want to treat the sensor like a micro-service.

* You want to use Golang.

### You don't support too many sensors...

No problem, you can help change that by contributing via [pull requeests](https://github.com/bartmika/treader-server/pulls) code for any sensors you think this project should have. If you are looking for a specific sensor, create a [request issue](https://github.com/bartmika/treader-server/issues).

## Installation

Install the application.

```
go install github.com/bartmika/treader-server
```

## Usage
Run our server continuously in the foreground:

```bash
$GOBIN/treader-server serve -f="/dev/cu.usbmodem14401" -s="SPARKFUN-DEV-13956"
```

If you see a message saying ``gRPC server is running.`` then the application has been successfully started.

The sub-command details are as follows:

```text
Run the gRPC server to allow other services to access the time-series data reader server

Usage:
  treader-server serve [flags]

Flags:
  -f, --arduino_path string     The location of the connected arduino device on your computer. (default "/dev/cu.usbmodem14201")
  -s, --arduino_shield string   The shield hardware attached to the arduino. (default "SPARKFUN-DEV-13956")
  -h, --help                    help for serve
  -p, --port int                The port to run this server on (default 50052)
```

Example output of successful operation:

```

```

## License

[**BSD 3-Clause License**](LICENSE) © Bartlomiej Mika
