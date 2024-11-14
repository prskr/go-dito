# Request handler

Request handlers are used in the plain HTTP server and will also be used in the GraphQL configuration.
Handlers are basically function calls that are modifying the HTTP response.
Similar to many programming languages, handlers can be polymorph, i.e. there are possibly multiple overload of the same hander.

## JSON handler

The `json(...)` handler is intended for small inline JSON responses.
It validates whether the passed string is a proper JSON string, writes the string to the response body and sets the `Content-Type` header to `application/json`.

There are multiple overloads:

1. `json(jsonResponse string)`
1. `json(statusCode int, jsonResponse string)`

In case there's no explicit status code given, dito will fallback to HTTP OK (200).


## File handler

The `file(...)` handler is more generic than the `json(...)` handler because it reads arbitrary files from the file system to the HTTP response.

The following overloads are available:

1. `file(filePath string)`
1. `file(filePath string, contentType string)`
1. `file(statusCode int, filePath string, contentType string)`

In case there's no explicit status code given, dito will fallback to HTTP OK (200).
If no content type is specified `go-dito` will try to infer the content type based on some heuristics of the Go standard library.
