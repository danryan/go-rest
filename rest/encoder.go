package rest

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

var encoders = make(map[string]EncoderFunc)

type Encoder interface {
	Encode(v interface{}) error
}

type EncoderFunc func(w io.Writer) Encoder

func (c *Client) Encoder(w io.Writer) (Encoder, error) {
	if c.EncoderFunc != nil {
		return c.EncoderFunc(w), nil
	}

	format, err := mediaTypeFormat(clientMediaType(c))
	if err != nil {
		return nil, err
	}

	if enc, ok := encoders[format]; ok {
		return enc(w), nil
	}

	return nil, fmt.Errorf("Could not determine an available decoder. Please specify a client.EncoderFunc.")
}

func (c *Client) Encode(v interface{}, w io.Writer) error {
	if v == nil {
		return fmt.Errorf("Input is nil, nothing to encode.")
	}

	enc, err := c.Encoder(w)
	if err != nil {
		return err
	}

	return enc.Encode(v)
}

func RegisterEncoder(format string, fn EncoderFunc) {
	encoders[format] = fn
}

func init() {
	RegisterEncoder("json", func(w io.Writer) Encoder {
		return json.NewEncoder(w)
	})

	RegisterEncoder("xml", func(w io.Writer) Encoder {
		return xml.NewEncoder(w)
	})
}
