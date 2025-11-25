# Service Example Go

This is a toy project for me to test and explore.

## Hexagonal architecture

The project show the advantages of using hexagonal architecture.
I always liked the use of the adapter pattern to decouple storage implementation,
well this is basically the same thing but for the driving side. This decouples the
code for exposing the service through REST / gRPC / ... and somewhat simplifies
testing as well.

## Why HUMA?

- Generates OpenAPI from the actual code and embeds a swagger for it
- Documentation is embedded with the code making required changes obvious
- Makes use of generics to use idiomatic functions as HTTP handlers.
- Handles errors and marshalling consistently

## Logging

The server is configured with some basic logging features:

- structured logging with `log/slog`
- access logs inspired by the [Common Log Format]
- panic recovery & logging
- service error logs
- X-Request-ID for basic correlation

[Common Log Format]: https://en.wikipedia.org/wiki/Common_Log_Format

## Metrics

The API expose those metrics:

```plaintext
build_info{goversion,title,version,revision,created}
http_requests_in_flight{method,path}
http_request_duration_seconds_bucket{method,path,status,le}
http_request_duration_seconds_sum{method,path,status}
http_request_duration_seconds_count{method,path,status}
http_requests_total{method,path,status}
process_*
```

This allow for request rate, error rate, concurrency, latency percentiles, averages...

## CI / CD

### Testing & Linting

There is basic github actions for triggering test and lint runs. The linter
(golangci) configuration is pretty aggressive but in my experience, most of the
time it ended up to be for the better. Almost never been in the way.

### Releasing

Making use of Goreleaser to streamline the process. The configuration is pretty
basic but it probably fits a good amount of real world use cases.

It builds the project according to an OS / architecture matrix, builds a docker
image, pushes it to a registry & creates a changelog.

### Updating

Basic Renovate configuration, making use of the Mend-hosted bot. It is set to
update all Go modules.

## Devcontainers

The project embeds its development environment through Development Containers.

## Project structure

```plaintext
.
|-- domain         The actual domain application code
|-- domaintest     Primitives to ease testing of domain components
|-- stores         Implementations of storage interfaces
|-- restapi        Registration of HTTP handlers for exposing a REST API
|-- router         Application agnostic routing helpers
|-- cli            Command-line facing objects & their options
`-- dist           For Goreleaser to use
```
