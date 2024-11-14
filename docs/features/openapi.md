# OpenAPI

When already working with OpenAPI (or [TypeSpec](https://typespec.io/)) you can use the schema to configure the mock server.
The server will generate the responses either based on the examples provided in the schema or the response schema.

If desired, you can also provide multiple examples and configure rules when to return which example.
To keep the schema readable and compatible with other tools, `dito` uses [specification extensions](https://swagger.io/specification/v3/#specification-extensions) to add conditions to your examples.

A simple example would be:

```yml
# ...
responses:
  "200":
    description: Successful operation
    content:
      application/json:
        schema:
          $ref: "#/components/schemas/Pet"
        examples:
          ted:
            value: |
              {"id":12, "name": "ted"}
            x-dito/when: 'http.JsonPath("$.name", "doggie")'
# ...
```

The conditions are expressed in a simple domain specific language (DSL) that allows you to match against various request properties.

Besides of the response body generation, `go-dito` also validates request and response against the schema and returns an error if the request or response doesn't match the schema.
This guarantees that the mock server behaves like the real API.
