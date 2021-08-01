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

```protobuf
service Telemetry {
    rpc GetTimeSeriesData (google.protobuf.Empty) returns (stream TelemetryDatum) {}
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

### Anyone using this?

The following project is currently using this project:

* [tpoller-server](https://github.com/bartmika/tpoller-server) - *Application written in Go which polls Time-series data at specific intervals and saves to persistent storage*
* [tstorage-server](https://github.com/bartmika/tstorage-server) - *Fast time-series data storage server accessible over gRPC*
* [sparkfunweathershield-arduino](https://github.com/bartmika/sparkfunweathershield-arduino) - *Application implemented in C++ which measures time-series data from the SparkFun Weather Shield board and provides serial usb interface for other devices*

If you want to add your own, feel free to make a [PR](https://github.com/bartmika/treader-server/pulls).

## Installation

Install the application.

```bash
go install github.com/bartmika/treader-server@latest
```

## Usage
Run our server continuously in the foreground:

```bash
$GOBIN/treader-server serve -f="/dev/cu.usbmodem14401" -s="SPARKFUN-DEV-13956" -p=50052
```

If your console output looks as the following then the application has been successfully started. You are ready to use the service!

```text
2021/07/15 22:00:16 READER: Attempting to connect Arduino device...
2021/07/15 22:00:16 READER: Waiting for Arduino external sensors to warm up
2021/07/15 22:00:26 gRPC server is running.
```

The sub-command details are as follows:

```text
Run the gRPC server to allow other services to access the time-series data reader

Usage:
  treader-server serve [flags]

Flags:
  -f, --arduino_path string     The location of the connected arduino device on your computer. (default "/dev/cu.usbmodem14201")
  -s, --arduino_shield string   The shield hardware attached to the arduino. (default "SPARKFUN-DEV-13956")
  -h, --help                    help for serve
  -p, --port int                The port to run this server on (default 50052)
```

Please note for `arduino_shield` the only supported values are:

* ``SPARKFUN-DEV-13956`` for [sparkfunweathershield-arduino](https://github.com/bartmika/sparkfunweathershield-arduino) code repo.

## Contributing
### Development
If you'd like to setup the project for development. Here are the installation steps:

1. Go to your development folder.

    ```bash
    cd ~/go/src/github.com/bartmika
    ```

2. Clone the repository.

    ```bash
    git clone https://github.com/bartmika/treader-server.git
    cd treader-server
    ```

3. Install the package dependencies

    ```bash
    go mod tidy
    ```

4. In your **terminal**, make sure we export our path (if you haven’t done this before) by writing the following:

    ```bash
    export PATH="$PATH:$(go env GOPATH)/bin"
    ```

5. Run the following to generate our new gRPC interface. Please note in your development, if you make any changes to the gRPC service definition then you'll need to rerun the following:

    ```bash
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/telemetry.proto
    ```

6. You are now ready to start the server and begin contributing!

    ```bash
    go run main.go serve -f="/dev/cu.usbmodem14401" -s="SPARKFUN-DEV-13956" -p=50052
    ```

### Quality Assurance

Found a bug? Need Help? Please create an [issue](https://github.com/bartmika/treader-server/issues).

## License

[**ISC License**](LICENSE) © Bartlomiej Mika
