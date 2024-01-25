# qs [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov]

This package is a fork of [github.com/pasztorpisti/qs](https://github.com/pasztorpisti/qs/) and contains 
several additional fixes and features that are not present in the upstream.

The `qs` package can marshal and unmarshal structs into/from url query strings.
The interface of `qs` is very similar to that of some standard marshaler
packages like `encoding/json`, `encoding/xml`.

Note that html forms are often `POST`-ed in the HTTP request body in the same
format as query strings (with an encoding called `application/x-www-form-urlencoded`)
so this package can be used for that as well.

# Quick Intro

The go standard library can convert only between the (query) string and the
standard [`url.Values` data type](https://golang.org/pkg/net/url/#Values)
(which is a `map[string][]string`).
This `qs` package adds struct marshaling and unmarshaling to your arsenal:

```
                  +--------------+
+---------------->| query string +------------------+
|                 +---------+----+                  |
|                      ^    |                       |
|    url.Values.Encode |    | url.ParseQuery        |
|                      |    v                       |
|                 +----+---------+                  |
|                 |  url.Values  |                  |
|                 +---------+----+                  |
|                      ^    |                       |
|     qs.MarshalValues |    | qs.UnmarshalValues    |
|                      |    v                       |
|                 +----+---------+                  |
+-----------------+    struct    |<-----------------+
    qs.Marshal    +--------------+    qs.Unmarshal
```

Example:
```go
package main

import "fmt"
import "github.com/pamburus/qs"

type Query struct {
	Search     string
	Page       int
	PageSize   int
	Categories []string `qs:"category"`
}

func main() {
	queryStr, err := qs.Marshal(&Query{
		Search:     "my search",
		Page:       2,
		PageSize:   50,
		Categories: []string{"c1", "c2"},
	})
	fmt.Println("Marshal-Result:", queryStr, err)

	var q Query
	err = qs.Unmarshal(&q, queryStr)
	fmt.Println("Unmarshal-Result:", q, err)

	// Output:
	// Marshal-Result: category=c1&category=c2&page=2&page_size=50&search=my+search <nil>
	// Unmarshal-Result: {my search 2 50 [c1 c2]} <nil>
}
```

# Features

- Support for primitive types (`bool`, `int`, etc...), pointers, slices, arrays,
  maps, structs, `time.Time` and `url.URL`.
- A custom type can implement the `MarshalQS` and/or `UnmarshalQS` interfaces
  to [handle its own marshaling/unmarshaling](https://godoc.org/github.com/pamburus/qs/#example-package--SelfMarshalingType).
- The marshaler and unmarshaler are modular and
  [can be extended to support new types](https://godoc.org/github.com/pamburus/qs/#example-package--CustomMarshalerFactory).
  This makes it possible to do several tricks. One of them is being able to
  override existing type marshalers (e.g.: the `[]byte` array marshaler).
- It can tell whether a type is marshallable before actually marshaling an
  object of that type. Most marshalers (including the standard `encoding/json`)
  can't do this! E.g.: If you have an empty slice that has a non-marshallable
  item type (e.g.: function) then `encoding/json` and many other marshalers
  happily marshal it as an empty slice - they return with an error only if the
  given slice has at least one item and they realise that it can't be marshaled.
- You can create custom marshaler objects that define different defaults for
  the marshaling process:
  - A struct-to-query_string name transformer func that is used when the struct
    field tag doesn't set a custom name for the field. The default function
    converts CamelCase go struct field names to snake_case which is the standard
    in case of query strings.
  - When a struct field tag specifies none of the  `keepempty` and `omitempty`
    options the marshaler uses `keepempty` by default. By creating a custom
    marshaler you can
    [change the default to `omitempty`](https://godoc.org/github.com/pamburus/qs/#example-package--DefaultOmitEmpty).
  - When a struct field tag doesn't specify any of the `opt`, `nil`, `req`
    options the unmarshaler uses `opt` by default. By creating a custom
    unmarshaler you can change this default.
- A struct field tag can be used to:
  - Exclude a field from marshaling/unmarshaling by specifying `-` as the
    field name (`qs:"-"`).
  - Set custom name for the field in the marshaled query string.
  - Set one of the `keepempty`, `omitempty` options for marshaling.
  - Set one of the `opt`, `nil`, `req` options for unmarshaling.

# Detailed Documentation

The [godoc of the qs package](https://godoc.org/github.com/pamburus/qs/)
contains more detailed documentation with working examples.


[doc-img]: https://pkg.go.dev/badge/github.com/pamburus/qs
[doc]: https://pkg.go.dev/github.com/pamburus/qs
[ci-img]: https://github.com/pamburus/qs/actions/workflows/ci.yml/badge.svg
[ci]: https://github.com/pamburus/qs/actions/workflows/ci.yml
[cov-img]: https://codecov.io/gh/pamburus/qs/graph/badge.svg?token=WBQLIRQGBO
[cov]: https://codecov.io/gh/pamburus/qs