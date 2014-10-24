# go-rest

A simple toolkit for building more useful HTTP clients.

Tested on all Go versions 1.0 and higher.

## Getting started

```go
package main

import "fmt"
import "github.com/danryan/go-rest/rest"

type githubUser struct {
  Login string `json:"login"`
  Name  string `json:"name"`
}

func main() {
  c, _ := rest.New("https://api.github.com", nil)
  c.Header.Set("Content-Type", "application/json")
  c.Header.Set("Accept", "application/json")

  user := &githubUser{}

  _, err := c.Get("users/danryan", user)
  if err != nil {
    panic(err)
  }

  fmt.Printf("%s's name is %s\n", user.Login, user.Name)
}

```

## Resources

* [API documentation](http://godoc.org/github.com/danryan/go-rest/rest)
* [Bugs, questions, and feature requests](https://github.com/danryan/go-rest/issues)

## Is it any good?

[Possibly.](http://news.ycombinator.com/item?id=3067434)

## License

This library is distributed under the MIT License, a copy of which can be found in the [LICENSE](LICENSE) file.
