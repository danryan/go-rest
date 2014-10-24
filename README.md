# go-rest

PagerDuty API client in Go.

Tested on all Go versions 1.0 and higher.

## Getting started

```go
package main

import "fmt"
import "github.com/danryan/go-rest/rest"

// Client is an API client
type Client struct {
  client    *rest.Client
}

type githubUser struct {
  Login string `json:"login"`
  Name string `json:"name"`
}

func main() {
  c := rest.New("https://api.github.com")
  c.Header.Set("Content-Type", "application/json")
  c.Header.Set("Accept", "application/json")

  user := &githubUser{}

  res, err := c.Get("users/danryan", user)
  if err != nil {
    panic(err)
  }

  fmt.Printf("%s's name is %s\n", u.Login, u.Name)
}
```

## Resources

* [API documentation](http://godoc.org/github.com/danryan/go-rest)
* [Bugs, questions, and feature requests](https://github.com/danryan/go-rest/issues)

## Is it any good?

[Possibly.](http://news.ycombinator.com/item?id=3067434)

## License

This library is distributed under the MIT License, a copy of which can be found in the [LICENSE](LICENSE) file.
