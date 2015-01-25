package main

import "fmt"
import "github.com/danryan/go-rest/rest"

// Client is an API client
type Client struct {
	client *rest.Client
}

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
