# Plain HTTP server

The plain HTTP server uses the same DSL as in the OpenAPI support to match against requests and extends it with the possibility to configure a response behavior.
It's the simplest way to configure a mock server and is perfect for small APIs or to get started quickly.

The downside is that the `go-dito` configuration is disconnected from your API schema and you have to maintain the configuration separately.
Also it doesn't support validation of the request/response bodies.
This might change in the future - at least for the request body - but it would require a user provided schema definition for the request body.
