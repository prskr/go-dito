# Plain HTTP configuration

As already mentioned, the plain HTTP server configuration is the easiest way to get started.
To set up a simple HTTP server stub, you can start with the following configuration:

```pkl linenums="1"
amends "https://raw.githubusercontent.com/prskr/go-dito/refs/heads/main/assets/AppConfig.pkl"

domains {
    ["localhost:3498"] = new PlainRuleSpec {
      rules = Set(
        #"http.Method("GET") -> http.Path("/api/v1/account/42") => File("testdata/sample.json", "application/json")"#,
        #"http.Method("POST") -> http.Path("/api/v1/account/42/withdraw") => Json(`{"name":"Ted.Tester"}`)"#
      )
    }
}
```

In the above snippet are two rules configured:

1. Respond to `GET /api/v1/account/42` with a file located at `testdata/sample.json` and set the content type header to `application/json`.
1. Respond to `POST /api/v1/account/42/withdraw` with a JSON body `#!json {"name":"Ted.Tester"}` - this also sets the content type to `application/json`

For further details on the DSL see [matchers](../dsl/matchers.md) and [handlers](../dsl/handlers.md)
