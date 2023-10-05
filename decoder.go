package jsonreplace

import (
	"encoding/json"
	"io"
)

// Unmarshal is like json.Unmarshal, but execute Unmarshal after replacing JSON objects.
func Unmarshal(data []byte, v interface{}, replacer Replacer) error {
	if replacer == nil {
		replacer = DefaultReplaceMux
	}
	if replacer == VoidReplacer {
		return json.Unmarshal(data, v)
	}
	bs, err := replacer.ReplaceJSON(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}

// Decoder is a JSON decoder with replaces JSON objects.
type Decoder struct {
	*json.Decoder
	replacer Replacer
}

// NewDecoder returns a new Decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		Decoder: json.NewDecoder(r),
	}
}

// Decode reads the next JSON-encoded value from its input and stores it in the value pointed to by v,
func (d *Decoder) Decode(v interface{}) error {
	replacer := d.replacer
	if replacer == nil {
		replacer = DefaultReplaceMux
	}
	if replacer == VoidReplacer {
		return d.Decoder.Decode(v)
	}
	var bs json.RawMessage
	if err := d.Decoder.Decode(&bs); err != nil {
		return err
	}
	bs, err := replacer.ReplaceJSON(bs)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}

// SetReplacer sets the Replacer.
func (d *Decoder) SetReplacer(replacer Replacer) {
	d.replacer = replacer
}
