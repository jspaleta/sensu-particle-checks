# Sensu Go Particle Checks
TravisCI: [![TravisCI Build Status](https://travis-ci.org/jspaleta/sensu-particle-checks.svg?branch=master)](https://travis-ci.org/jspaleta/sensu-particle-checks)

Collection of Sensu Agent check commands to interact with Particle Cloud

## Installation

Download the latest version of the sensu-particle-checks from [releases][1],
or create an executable script from this source.

### Build from source
From the local path of the sensu-particle-checks repository:

#### Device Variable Check
```
go build -o /usr/local/bin/particle_variable_check ./particle_variable_check/main.go
```

#### Device Ping Check
```
go build -o /usr/local/bin/particle_ping_check ./particle_ping_check/main.go
```

## Configuration

Example Sensu Go definition:

```json
{
    "api_version": "core/v2",
    "type": "check",
    "metadata": {
        "namespace": "default",
        "name": "particle"
    },
    "spec": {
        "...": "..."
    }
}
```

## Usage Examples

### Device Variable Check

Help:

### Device Ping Check

## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/jspaleta/sensu-particle-checks/releases
