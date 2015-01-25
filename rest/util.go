package rest

import (
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"reflect"
	"regexp"

	"github.com/google/go-querystring/query"
)

// Bool is a helper routine that allocates a new bool value to store v and
// returns a pointer to it.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// Int is a helper routine that allocates a new int32 value to store v and
// returns a pointer to it, but unlike Int32 its argument value is an int.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// String is a helper routine that allocates a new string value to store v
// and returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}

// AddOptions adds the parameters in opt as URL query parameters to s.
// opt must be a struct whose fields may contain "url" tags.
func AddOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()

	return u.String(), nil
}

func clientMediaType(c *Client) string {
	if m := c.Header.Get("Accept"); m != "" {
		return m
	}
	if m := c.Header.Get("Content-Type"); m != "" {
		return m
	}

	return ""
}

func responseMediaType(r *http.Response) string {
	if m := r.Header.Get("Content-Type"); m != "" {
		return m
	}

	return ""
}

func mediaTypeFormat(mtype string) (string, error) {
	mt, _, err := mime.ParseMediaType(mtype)
	if err != nil {
		return mt, err
	}

	known := []string{"json", "xml"}
	for _, f := range known {
		if match, _ := regexp.MatchString(f, mt); match {
			return f, nil
		}
	}

	return "", fmt.Errorf(`Could not determine the media type from %s.`, mtype)
}

// BasicAuth takes a username and password string and returns a base64-encoded string
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
	// return base64.URLEncoding.EncodeToString([]byte(auth))
}
