[![Go build](https://github.com/Netcracker/qubership-core-lib-go-actuator-common/actions/workflows/go-build.yml/badge.svg)](https://github.com/Netcracker/qubership-core-lib-go-actuator-common/actions/workflows/go-build.yml)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?metric=coverage&project=Netcracker_qubership-core-lib-go-actuator-common)](https://sonarcloud.io/summary/overall?id=Netcracker_qubership-core-lib-go-actuator-common)
[![duplicated_lines_density](https://sonarcloud.io/api/project_badges/measure?metric=duplicated_lines_density&project=Netcracker_qubership-core-lib-go-actuator-common)](https://sonarcloud.io/summary/overall?id=Netcracker_qubership-core-lib-go-actuator-common)
[![vulnerabilities](https://sonarcloud.io/api/project_badges/measure?metric=vulnerabilities&project=Netcracker_qubership-core-lib-go-actuator-common)](https://sonarcloud.io/summary/overall?id=Netcracker_qubership-core-lib-go-actuator-common)
[![bugs](https://sonarcloud.io/api/project_badges/measure?metric=bugs&project=Netcracker_qubership-core-lib-go-actuator-common)](https://sonarcloud.io/summary/overall?id=Netcracker_qubership-core-lib-go-actuator-common)
[![code_smells](https://sonarcloud.io/api/project_badges/measure?metric=code_smells&project=Netcracker_qubership-core-lib-go-actuator-common)](https://sonarcloud.io/summary/overall?id=Netcracker_qubership-core-lib-go-actuator-common)

# Go-actuator-common
* [Install](#install)
* [Health-core](#health-core)
* [Api-version](#api-version)
* [Configloader-property-utils](#config-property-utils)
* [Monitoring-core (internal)](#monitoring-core)
* [Tracer-core (internal)](#tracer-core)
* [Log levels (internal)](#log-levels)

> **_NOTE:_** Internal modules are not public modules and are intended solely for internal use. Any methods can be changed without backward compatibility.

# Install

To get `go-actuator-common` use
```go
 go get github.com/netcracker/qubership-core-lib-go-actuator-common/v2@<latest released version>
```

# Health-core

Health core paclage contains functionality which is responsible for processing health indicators and calculate health status.

### How to use

`First of all` you have to create a `healthService` instance. You can use `healthServiceImp` implementation or create your own.
```go
import (
"github.com/netcracker/qubership-core-lib-go-actuator-common/v2/health"
)

healthService := health.NewHealthService(<HealthConfiguration>)
```
`HealthConfiguration` is an `optional` struct, that contains some configurations attributes, for example `HealthCheckInterval` a value which indicates a delay between checks. 

After that, you have to add `health checks`. For that you should implement `Check` function and provide 
using `HealthService#AddCheck` function.

`Secondly`, you need to start `healthService`.
```go
healthService.Start()
```
When you call method `start` healthService will start to calculate health status by provided `checks`.

To get result, you have to call the method:
```go
data := healthService.GetHealth()
```

You can notice that we don't provide any rest handler. So, for using it in a rest microservice 
you need to create your own handler in which you will call `GetHealth()` method.

# Api-version

apiversion package contains functionality which is responsible for processing api versioning information.

### How to use

`First of all` you have to create a `apiVersionService` instance. You can use `apiVersionImpl` implementation.
```go
import (
    "github.com/netcracker/qubership-core-lib-go-actuator-common/v2/apiversion"
)

config := apiversion.ApiVersionConfig{
    PathToApiVersionInfoFile: "../config/api-version-info.json",
}

apiVersionService, _ := apiversion.NewApiVersionService(config)
```
Use `ApiVersionConfig` for setting path to json file with api configuration.
Default value for PathToConfigFile is `./api-version-info.json`. You can set different path to json file with config.

Also you can use property `"api-version.path"` to set path to file.

To get result, you have to call the `GetApiVersion(ctx context.Context)` method:
```go
data, err := apiVersionService.GetApiVersion(ctx)
```
This method return `*ApiVersionResponse`

You can notice that we don't provide any rest handler. So, for using it in a rest microservice
you need to create your own handler in which you will call `GetApiVersion(ctx context.Context)` method.

Example of file:
```json
{
  "specs": [
    {
      "specRootUrl": "/api",
      "major": 3,
      "minor": 44,
      "supportedMajors": [
        2,
        3
      ]
    },
    {
      "specRootUrl": "/api/bluegreen",
      "major": 1,
      "minor": 10,
      "supportedMajors": [
        1
      ]
    }
  ]
}
```
All fields (specRootUrl, major, minor, supportedMajors) are required and cannot be empty.

# Config-property-utils

This package contains some common preconfigured convenient methods that allow get property value from configloader
property sources.  
In order to use it, you should import ```clpropertyutils``` package and call some util API method. For example:

```go
import (
    clpropertyutils "github.com/netcracker/qubership-core-lib-go-actuator-common/v2/configloader-property-utils"
)

if config.ReadBufferSize <= 0 {
    config.ReadBufferSize = clpropertyutils.GetHttpBufferHeaderMaxSizeBytes()
    logger.Info("HTTP buffer header Max size has been set to %d bytes", config.ReadBufferSize)
}
```

## Provided methods

| Method                                | Description                                                    | Default value |
|---------------------------------------|----------------------------------------------------------------|---------------|
| GetHttpBufferHeaderMaxSizeBytes():int | Max http header buffer size in KB that is used by http server. | 10240 (KB)    |

#  Monitoring-core
> **_NOTE:_** This is not public module and only intended for internal using, any methods can be changed without backward compatibility. 

`Monitoring-core` module represents metrics which we recommend and which are used in our libraries.   
Prometheus metrics:
 * request counter;
 * latency counter.

To start using them, you need   

Firstly, to call:
```go
import (
"github.com/netcracker/qubership-core-lib-go-actuator-common/v2/monitoring"
)

platformPrometheusMetrics, err := monitoring.RegisterPlatformPrometheusMetrics()
```

By calling this function you register metrics in the global prometheus instance and get them.


Secondly, you need to create a middleware where you will register request and increment metrics by using methods:
*  `IncRequestStatusCounter`
 * `ObserveRequestLatencyHistogram`.
 
This methods must be called after client's rest handler.

Thirdly,you need to add an endpoint `promhttp.Handler()` which will return collected metrics.

# Tracer-core
> **_NOTE:_** This is not public module and only intended for internal using, any methods can be changed without backward compatibility. 

Tracer helps to trace requests. `go-actuator-tracing` uses standard OpenTelemetry and provides zipkin opentelementary implementation.  
If you want to use a different tracer, you have to implement OpenTelemetryExporter.

To register zipkin tracer you need to create `zipkinTracer` and call `RegisterTracerProvider` method. There two ways to initiate `zipkinTracer` using environment parameters 
or explicitly pass ones. In order to init using env parameters you should pass them:

|Name|Description|Default|Allowed values|
|---|---|---|---|
|tracing.enabled  | Enable or disable tracing (to switch on/off without changing other params) | false | true/false|
|tracing.host     | Zipkin host server, without port and protocol | -- | any string, for example nc-diagnostic-agent
|tracing.sampler.const    | sampler always makes the same decision for all traces. It either samples all traces (value=1) or none of them (value=0). | 1 | 0 or 1 
|microservice.name    | microservice name | -- | any string, for example tenant-manager

```go
import (
	"go.opentelemetry.io/otel"
)

zipkinTracer := NewZipkinTracer()
registered, err := zipkinTracer.RegisterTracerProvider()
err := zipkinTracer.RegisterTracerProvider()
```

Besides it, you can use a constructor with `ZipkinOptions` parameter and specify them yourself. For example:
```go
import (
	"go.opentelemetry.io/otel"
)

options := ZipkinOptions{ServiceName: "someService", TracingHost: "localhost", TracingSamplerConst: 1, TracingEnabled: true}
zipkinTracer := NewZipkinTracerWithOpts(options)
registered, err := zipkinTracer.RegisterTracerProvider()
err := zipkinTracer.RegisterTracerProvider()
```

After that you can use zipkin tracer as an opentelementary tracer provider.

# Log levels
> **_NOTE:_** This is not public module and only intended for internal using, any methods can be changed without backward compatibility. 

Log level service implements logic for getting log levels for all currently created loggers.

Usage example:
```go
import (
        "github.com/netcracker/qubership-core-lib-go-actuator-common/v2/loglevel"
)

logLevelService, _ = loglevel.NewLogLevelService()
data, err := logLevelService.GetLogLevels()
```
