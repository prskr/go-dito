# Configuration

`go-dito` can be configured via a configuration file in [PKL](https://pkl-lang.org/) format.
The schema of the configuration can be found in the [AppConfig.pkl](https://github.com/prskr/go-dito/blob/main/assets/AppConfig.pkl)

To override the default configuration you have to `amend` the configuration file:

```pkl
amends "https://raw.githubusercontent.com/prskr/go-dito/refs/heads/main/assets/AppConfig.pkl"

server {
  port = 8080
}
```

The `amends` statement is used to include the default configuration and override it with your own configuration.
In the example above it is using the latest version of the main branch, as soon as there are releases you should use the release tag instead of the branch name.

The advantage of PKL is, that the configuration has a schema and you can get autocompletion and validation in your [editor](https://pkl-lang.org/main/current/tools.html).

## Domains

`domains` is the central configuration section where you can configure multiple domains with different behaviors.
It is only possible to use **one kind** of behavior **per domain**.
This means, it is not possible to mix for instance OpenAPI and GraphQL on a single domain.
This might change in the future if necesssary but for now this keeps the complexity both in the configuration and in the implementation at bay.

## Server

The `server` section is where listening host and port are configured.
Furthermore there are some fine grained configuration options for the HTTP server such as `readHeaderTimeout`.

## Telemetry

The `telemetry` section is where things like logging is configured and also OpenTelemetry (OTeL) related settings will be located in this section.
