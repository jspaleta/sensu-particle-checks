# Sensu Go Particle Checks
TravisCI: [![TravisCI Build Status](https://travis-ci.com/jspaleta/sensu-particle-checks.svg?branch=master)](https://travis-ci.com/jspaleta/sensu-particle-checks)

Sensu Go Particle Checks is a collection of Sensu Go Agent [check commands](https://docs.sensu.io/sensu-go/5.0/reference/checks/#how-do-checks-work) that interact with [Particle Device Cloud API](https://docs.particle.io/reference/device-cloud/api/).  

These check commands will allow you to schedule interaction with the Particle Device Cloud using Sensu Go, and to produce Sensu alert and metrics based on on Particle Device Cloud information.

## Installation

Download the latest version of the sensu-particle-checks from [releases][1],
or create an executable script from this source.

### Build from source
Install the source into your `GOPATH/src` path. 

You can find your configured GOPATH using `go env`
```
$ go env
GOARCH="amd64"
...
GOPATH="/home/myuser/go"
...
```

The easiest way to do that is with `go get`
```
go get github.com/jspaleta/sensu-particle-checks
```

Then from the local path of the sensu-particle-checks repository:

#### Build device variable check
```
go build -o /usr/local/bin/particle_variable_check ./particle_variable_check/main.go
```

#### Build device ping check
```
go build -o /usr/local/bin/particle_ping_check ./particle_ping_check/main.go
```

## Configuration

Example Sensu Go check definitions:

### Particle variable metric check

```json
{
  "type": "CheckConfig",
  "api_version": "core/v2",
  "metadata": {
    "name": "particle_variable_metric",
    "namespace": "default"
  },
  "spec": {
    "command": "particle_variable_metric_check -v variable -m metricname.variable",
    "env_vars": [
      "PARTICLE_DEVICEID=123456789123456789",
      "PARTICLE_TOKEN=abcdefgABCDEF12345"
    ],
    "handlers": [
      "status"
    ],
    "interval": 60,
    "output_metric_format": "graphite_plaintext",
    "output_metric_handlers": [
      "influxdb"
    ],
    "publish": true,
    "subscriptions": [
      "particle"
    ],
    "timeout": 30
  }
}
```

### Particle ping check

```json
{
  "type": "CheckConfig",
  "api_version": "core/v2",
  "metadata": {
    "name": "particle_ping_check",
    "namespace": "default"
  },
  "spec": {
    "command": "particle_variable_metric_check",
    "env_vars": [
      "PARTICLE_DEVICEID=123456789123456789",
      "PARTICLE_TOKEN=abcdefgABCDEF12345"
    ],
    "handlers": [
      "status"
    ],
    "interval": 60,
    "publish": true,
    "subscriptions": [
      "particle"
    ],
    "timeout": 30
  }
}
```
**Security Note:** The Particle access token, deviceID and productID are treated as a security sensitive configuration options in these examples and are loaded into the handler config as env_vars instead of as command arguments. Command arguments are commonly readable from the process table by other unprivaledged users on a system (ex: `ps` and `top` commands), so it's a better practise to read in sensitive information via environment variables or configuration files as part of command execution. The command flags for these configuration options are provided as an override for testing purposes.


## Usage Examples

### Device Variable Check

Help: `particle_variable_metric_check --help`
```
Retrieve Particle string variable and output in graphite plaintext format

Usage:
  particle_variable_metric_check [flags]

Flags:
  -a, --access_token string   required Particle Access Token, defaults to PARTICLE_TOKEN env variable
  -d, --device string         required Particle DeviceID, defaults to PARTICLE_DEVICEID env variable
      --dryrun                dryrun to check inputs
  -h, --help                  help for particle_variable_metric_check
  -m, --metric string         optional metric name, if not set will be determined from hostname.variable
  -p, --product string        optional Particle ProductID, defaults to PARTICLE_PRODUCTID env variable
  -T, --timeout int           optional particle Metric Timestamp Timeout (seconds) (default 60)
  -t, --timestamp string      optional Particle Timestamp Variable, must hold string representation of Unix Epoch integer
  -v, --variable string       required Particle Variable name, must hold string value
      --verbose               enable verbose output
```

### Device Ping Check

Help: `particle_ping_check --help`
```
Ping Particle device and check to see if its online

Usage:
  particle_ping_check [flags]

Flags:
  -a, --access_token string   required Particle Access Token, defaults to PARTICLE_TOKEN env variable
  -d, --device string         required Particle DeviceID, defaults to PARTICLE_DEVICEID env variable
  -h, --help                  help for particle_ping_check
  -p, --product string        optional Particle ProductID, defaults to PARTICLE_PRODUCTID env variable
      --verbose               enable verbose output

```

## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/jspaleta/sensu-particle-checks/releases
