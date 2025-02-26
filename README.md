# go-dito

## Why?

'dito' is an API simulation tool that tries to cover the most crucial use cases when:

- you're trying to mock one or multiple APIs that you're working with
- the APIs you're working with are either
    - RESTful
    - GraphQL
    - HTTP

  based
- the mocking rules are - for now - static, meaning, you don't want to manipulate the state of your mocks and don't rely
  on dynamic response generation
- you probably want to trace the interaction with your mock server to not lose context

## Non-goals - for now

- state store to dynamically manipulate mocks
- dynamic rules to respond based on previous state manipulation or external factors (e.g. by calling a 3rd API)

## Roadmap

- [ ] gRPC support
- [ ] JS/TS script support for dynamic request matching and response generation

### dito DSL

The DSL has the following schema:

```
[matcher [ -> matcher ]] => response
```

`matchers` are combined in a logical *and* fashion.
The more `matchers` you're chaining, the more specific a rule becomes.

When dito is parsing the rules it compares your matchers and response providers based on their signatures.
The signature has the following schema

```
[module.]name([[type], type...])
```

meaning: every rule can have a module (like `http` or `graphql`), it has to have a name and zero to multiple parameters
that are compared based on their type (e.g. `string` or `int`).

### Matchers

For now dito has a rather limited set of matchers:

| Signature                           | Description                                                                                              | Example                                                            |
|-------------------------------------|----------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------|
| `http.Method(string)`               | Match the HTTP method                                                                                    | `http.Method("GET")`                                               |
| `http.HeaderPresent(string)`        | Checks whether the request has a certain header                                                          | `http.HeaderPresent("Authorization")`                              |
| `http.Header(string, string)`       | Checks whether the request has a certain header value                                                    | `http.Header("Content-Type", "application/json")`                  |
| `http.Path(string)`                 | Compares the request path with the given value (case sensitive)                                          | `http.Path("/health")`                                             |
| `http.PathPattern(string)`          | Matches the given regex against the request path                                                         | `http.PathPattern("/health")`                                      |
| `http.Query(string, string)`        | Checks whether the request has a certain query value                                                     | `http.Query("limit", "100")`                                       |
| `http.QueryPattern(string, string)` | Matches the given regex against a query value                                                            | `http.QueryPattern("limit", "100")`                                |
| `http.JSONPath(string, string)`     | Extracts a value based on the given JSON path from the request body and compares it with the given value | `http.JSONPath("$.some.path", "hello")`                            |
| `graphql.Query(string)`             | Match the given GraphQL query against the one in the request body                                        | `graphql.Query("query { allFilms { films { director title } } }")` |
| `graphql.QueryFromFile(string)`     | Reads a GraphQL query from a file and compares it with the one in the request body                       | `graphql.QueryFromFile("testdata/queries/simple.gql")`             |

### Response providers

Every rule **has** to have a response provider that specifies what to do with the incoming request.

For now dito ships with the following response providers:

| Signature           | Description                                                               | Example                            |
|---------------------|---------------------------------------------------------------------------|------------------------------------|
| `Status(int)`       | Return only an HTTP status code                                           | `Status(202)`                      |
| `JSON(string)`      | Return an inline specified JSON response                                  | ``JSON(`{"hello":"world"}`)``      |
| `JSON(int, string)` | Return an inline specified JSON response and specify the HTTP status code | ``JSON(202, `{"hello":"world"}`)`` |
| `File(string)`      | Return the content of the specified file                                  | `File("testdata/simple.json")`     |

## Configuration

The configuration is [Pkl](https://pkl-lang.org/index.html) based.
The [`config.pkl`](./config.pkl) in this repository illustrates most configuration options - especially those based on
the DSL.
The schema of the configuration is [`AppConfig.pkl`](./assets/AppConfig.pkl).
See also the [docs](https://pkl-lang.org/main/current/language-tutorial/02_filling_out_a_template.html#amending) on how
to amend a custom configuration.

## Alternatives

There are great alternatives that share at least some of the functionality of dito:

- [Microcks](https://microcks.io/)
- [smocker](https://smocker.dev/)
