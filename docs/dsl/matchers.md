# Request matchers

Request matchers are used in basically every kind of server configuration.
They are matching against properties of an incoming HTTP request and are helping `go-dito` to distinguish between multiple handlers.

Matchers can be 'chained' to build sophisticated routing rules.

Similar to many programming languages, matchers can be polymorph, i.e. there are possibly multiple overload of the same matcher.

## HTTP method

The `http.method(methodName string)` matcher - as the name already suggests - matches on the HTTP method.
It's a good start for a matcher chain because it's very cheap and can be used to filter out requests early.


## HTTP path

The `http.path(path string)` matcher matches on the request path.
It is an **exact** match, so it's not possible to use wildcards or regex here.
It is also **case sensitive**.
For this reason this matcher is also a lot cheaper than for instance [HTTP path pattern](#http-path-pattern) and should therefore be preferred if possible.

## HTTP path pattern

The `http.pathPattern(pattern string)` matcher is a more sophisticated version of the [HTTP path](#http-path) matcher.
It allows to use regex patterns to match on the request path.
The pattern is compiled to a regex and matched against the request path.
The underlying regex engine is the [Go standard library](https://pkg.go.dev/regexp) regex engine.

## HTTP header

The `http.header(key string, value string)` matcher matches on a specific header key and value.
It is also an **exact** match, so it's not possible to use wildcards or regex here.
For this reason this matcher is also considered a cheap matcher and can should be rather early in the matcher chain.

## HTTP header present

The `http.headerPresent(key string)` matcher is an easier variant of [HTTP header](#http-header).
It only matches whether a certain header is present, but it does not check the value of the header at all.

## HTTP query

The `http.query(key string, value string)` matcher checks whether a certain combination of query key and value are set for an incoming HTTP request.
The matcher is both **exact** and **case sensitive**.
A more generic alternative is [HTTP query pattern](#http-query-pattern) but due to the fact that it is based on regex patterns, it is also less performant.

## HTTP query pattern

The `http.queryPattern(key string, pattern string)` matcher is similar to the aforementioned [HTTP query](#http-query) matcher but uses a regex pattern instead of an exact string to match the query value.
