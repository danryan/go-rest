package rest

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

var decoders = make(map[string]DecoderFunc)

type Decoder interface {
	Decode(v interface{}) error
}

type DecoderFunc func(r io.Reader) Decoder

func (c *Client) Decoder(r io.Reader) (Decoder, error) {
	if c.DecoderFunc != nil {
		return c.DecoderFunc(r), nil
	}

	format, err := mediaTypeFormat(clientMediaType(c))
	if err != nil {
		return nil, err
	}

	if dec, ok := decoders[format]; ok {
		return dec(r), nil
	}

	return nil, fmt.Errorf("Could not determine an available decoder. Please specify a client.DecoderFunc.")
}

func (c *Client) Decode(v interface{}, r io.Reader) error {
	if v == nil {
		return nil
	}

	dec, err := c.Decoder(r)
	if err != nil {
		return err
	}

	return dec.Decode(v)
}

func RegisterDecoder(format string, fn DecoderFunc) {
	decoders[format] = fn
}

func init() {
	RegisterDecoder("json", func(r io.Reader) Decoder {
		return json.NewDecoder(r)
	})

	RegisterDecoder("xml", func(r io.Reader) Decoder {
		return xml.NewDecoder(r)
	})
}
