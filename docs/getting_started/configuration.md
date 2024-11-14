# Configuration

`go-dito` can be configured via a configuration file in [PKL](https://pkl-lang.org/) format although this might change in the future.
The schema of the configuration can be found in the [AppConfig.pkl](assets/AppConfig.pkl)

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

The central configuration is the `domains` section where you can configure multiple domains with different behavior.
