# Reverse-Geocoding-Service

[![Build Status](https://travis-ci.com/adaptant-labs/reverse-geocoding-service.svg?branch=master)](https://travis-ci.com/adaptant-labs/reverse-geocoding-service#)

A simple reverse geocoding Consul-enabled microservice wrapped around the [georeverse] package written in Go.

`reverse-geocoding-service` takes either a latitude and longitude as an encoded JSON pair or an IP address and returns
a corresponding [ISO 3166-1 alpha-2] 2-character country code.

[georeverse]: https://github.com/adaptant-labs/geo.git
[ISO 3166-1 alpha-2]: https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2

## Installation

```sh
go get github.com/adaptant-labs/reverse-geocoding-service
```

## Usage

On the server side, simply run the service directly:

```sh
$ reverse-geocoding-service
2019/06/30 21:28:42 Listening on :4041/georeverse...
```

The following options can be set at run-time:

```sh
$ reverse-geocoding-service --help
  NAME:
     Reverse Geocoding Service - Microservice for providing ISO 3166-1 country code lookup from a lat/lng pair
  
  USAGE:
     reverse-geocoding-service [global options] command [command options] [arguments...]
  
  VERSION:
     0.0.2
  
  AUTHOR:
     Adaptant Labs <labs@adaptant.io>
  
  COMMANDS:
       help, h  Shows a list of commands or help for one command
  
  GLOBAL OPTIONS:
     --port value          Port to bind to (default: 4041)
     --host value          Host address to bind to
     --data value          Polygon definition file to use (default: "./data/polygons.properties")
     --use-consul          Use Consul for Service Registration
     --consul-agent value  Consul Agent to connect to (default: "localhost:8500")
     --help, -h            show help
     --version, -v         print the version

  COPYRIGHT:
     Adaptant Solutions AG
```

From the client side, this can be tested by sending a JSON-encoded lat/lng pair:

```
$ curl -X POST -d '{ "lng": -89.234005, "lat": 41.645332 }' http://localhost:4041/georeverse
```

with the country shortcode returned in the POST response body:

```
{"country_code":"US"}
```

IP-based lookup follows a similar process:

```
$ curl -X POST http://localhost:4041/georeverse/8.8.8.8
{"country_code":"US"}
```
## Service Discovery

By default, the service will attempt to register itself with a local Consul server. The agent location can be tuned
with the `--consul-agent` command line flag, while service registration with Consul can be inhibited by setting
`--use-consul=false`. The service itself is registered under `reverse-geocoding`, as below:

```json
[
    {
        "ID": "fa89781d-2958-f1e6-16ff-6c7ad22f9ed0",
        "Node": "sgx-CELSIUS-W550power",
        "Address": "127.0.0.1",
        "Datacenter": "dc1",
        "TaggedAddresses": {
            "lan": "127.0.0.1",
            "wan": "127.0.0.1"
        },
        "NodeMeta": {
            "consul-network-segment": ""
        },
        "ServiceKind": "",
        "ServiceID": "reverse-geocoding",
        "ServiceName": "reverse-geocoding",
        "ServiceTags": [],
        "ServiceAddress": "",
        "ServiceWeights": {
            "Passing": 1,
            "Warning": 1
        },
        "ServiceMeta": {},
        "ServicePort": 4041,
        "ServiceEnableTagOverride": false,
        "ServiceProxyDestination": "",
        "ServiceProxy": {},
        "ServiceConnect": {},
        "CreateIndex": 99,
        "ModifyIndex": 99
    }
]
```

## Deployment

Docker images are provided under [adaptant/reverse-geocoding-service][docker] and can be run with and without
Consul-backed service registration. To run the service directly without Consul support:

```
$ docker run -d -p 4041:4041 adaptant/reverse-geocoding-service --use-consul=false
```

[docker]: https://hub.docker.com/r/adaptant/reverse-geocoding-service

While if being used together with a Consul Agent, the agent location can be specified as:

```
$ docker run -d -p 4041:4041 adaptant/reverse-geocoding-service --consul-agent=<agent ip>:<port>
```

## Datasets

This service uses the world borders dataset from [thematicmapping.org]. The data is stored in a properties file within
the project that maps country codes to Well Known Text format polygons and multi-polygons. The dataset was converted to
Well Known Text format by [GIS Stack Exchange user, elrobis](http://gis.stackexchange.com/a/17441).

This service further uses the [GeoLite2 Database](https://dev.maxmind.com/geoip/geoip2/geolite2/) and Contents, Copyright (c) 2019 MaxMind, Inc.

[thematicmapping.org]: http://thematicmapping.org/downloads/world_borders.php

## Features and bugs

Please file feature requests and bugs at the [issue tracker][tracker].

[tracker]: https://github.com/adaptant-labs/reverse-geocoding-service/issues

## License

`reverse-geocoding-service` itself is licensed under the terms of the Apache 2.0 license, the full version of which can
be found in the [LICENSE](LICENSE) file included in the distribution.

The country polygon dataset included in `data/polygons.properties` is made available under [CC BY-SA 3.0] in
accordance with the license of the [reverse-country-code] (since re-licensed under the Apache 2.0 license) project under
which it was originally released.

The GeoLite2 database is licensed under [CC BY-SA 4.0], and further incorporates [GeoNames](http://www.geonames.org)
geographical data, which is made available under [CC BY 3.0 US].

[CC BY 3.0 US]: http://www.creativecommons.org/licenses/by/3.0/us/
[CC BY-SA 3.0]: http://creativecommons.org/licenses/by-sa/3.0/
[CC BY-SA 4.0]: http://creativecommons.org/licenses/by-sa/4.0/
[reverse-country-code]: https://github.com/bencampion/reverse-country-code
