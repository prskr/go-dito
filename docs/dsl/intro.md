# DSL introduction

`go-dito` comes with a custom DSL to configure routing rules in a very simple yet powerful way.

Generally, the DSL itself - meaning names of matchers & providers are **not** case sensitive.
So when writing rules you can use `http.method("GET")` as well as `HTTP.Method("GET")` or `HTTP.method("GET")`.

"Arguments" of "functions" in the DSL **are potentially** case sensitive depending on the individual use case.
