# telemetry-server

## Overview

The purpose of this server is to provide a remote procedure call (gRPC) interface over an external Arduino device(with a shield); more specifically, this application implements the [`Telemetry`](https://github.com/bartmika/tpoller-server/blob/master/proto/telemetry.proto) gRPC service definition so as long as you implement that protocol then you can access the time-series data.

Currently this server only supports the following shield:

* [SparkFun Weather Shield (DEV-13956)](https://github.com/sparkfun/Weather_Shield).

## Prerequisites

You must have the following installed before proceeding. If you are missing any one of these then you cannot begin.

* ``Go 1.16.3``

## Installation

1. Please visit the [sparkfunweathershield-arduino](https://github.com/bartmika/sparkfunweathershield-arduino) repository and setup the external device and connect it to your development machine.

2. Please find out what USB port your external device is connected on. Note: please replace ``/dev/cu.usbmodem14201`` with the value on your machine, a Raspberry Pi would most likely have the value ``/dev/ttyACM0``.

3. Download the source code, build and install the application.

    ```
    GO111MODULE=on go get -u github.com/bartmika/telemetry-server
    ```

## Usage
Run our application.

```bash
$GOBIN/telemetry-server serve -f="/dev/cu.usbmodem14401" -s="SPARKFUN-DEV-13956"
```

If you see a message saying ``gRPC server is running.`` then the application has been successfully started.

## How does it work?
This device runs continuously waiting for you to pull data from it.

When you pull data, it will send you a JSON formatted object with all the time series data.

To pull data, you must first connect to the **Arduino device** with a USB cable.

Once connected, you use **serial usb communication** to read data from the device and write commands to the device.

Once your device recieves the JSON data, you do what you want with the data.

## Why did you choose Arduino?
The Arduino platform has a wonderful ecosystem of open-source hardware with libraries. Our goal is to take advantage of the libraries the hardware manufacturers wrote and not worry about the complicated implementation details.

## How does the data output look like?
When the device is ready to be used, you will see this output:

```json
{"status":"READY","runtime":2,"id":1,"sensors":["humidity","temperature","pressure","illuminance","soil"]}
```

When you poll the device for data, you will see this output:

```json
{"status":"RUNNING","runtime":24771,"id":2,"humidity":{"value":47.92456,"unit":"%","status":1,"error":""},"temperature_primary":{"value":80.47031,"unit":"F","status":1,"error":""},"pressure":{"value":0,"unit":"Pa","status":1,"error":""},"temperature_secondary":{"value":78.2375,"unit":"F","status":1,"error":""},"altitude":{"value":80440.25,"unit":"ft","status":1,"error":""},"illuminance":{"value":0.040305,"unit":"V","status":1,"error":""}}
```

## Why should I use it?
This code is a easy to connect and read realtime time-series data using any language that supports serial communication over USB.

## License

This application is licensed under the **BSD 3-Clause License**. See [LICENSE](LICENSE) for more information.
