package jsonreplace

import (
	"encoding/json"
	"io"
)

// Marshal returns the JSON encoding of v, replacing JSON objects as json.Marshal does.
func Marshal(v interface{}, replacer Replacer) ([]byte, error) {
	if replacer == nil {
		replacer = DefaultReplaceMux
	}
	bs, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if replacer == VoidReplacer {
		return bs, nil
	}
	return replacer.ReplaceJSON(bs)
}

// MarshalIndent is like Marshal but applies Indent to format the output.
func MarshalIndent(v interface{}, replacer Replacer, prefix, indent string) ([]byte, error) {
	if replacer == nil {
		replacer = DefaultReplaceMux
	}
	bs, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if replacer == VoidReplacer {
		return bs, nil
	}
	bs, err = replacer.ReplaceJSON(bs)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(json.RawMessage(bs), prefix, indent)
}

// Encoder is a JSON encoder with replaces JSON objects.
type Encoder struct {
	*json.Encoder
	replacer Replacer
}

// NewEncoder returns a new Encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		Encoder: json.NewEncoder(w),
	}
}

// Encode writes the JSON encoding of v to the stream, replacing JSON objects
// as json.Encoder does.
// if not SetReplacer, use DefaultReplaceMux
func (e *Encoder) Encode(v interface{}) error {
	replacer := e.replacer
	if replacer == nil {
		replacer = DefaultReplaceMux
	}
	if replacer == VoidReplacer {
		return e.Encoder.Encode(v)
	}
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	raw := json.RawMessage(bs)

	replaced, err := replacer.ReplaceJSON(raw)
	if err != nil {
		return err
	}
	return e.Encoder.Encode(replaced)
}

// SetReplacer sets the Replacer to use.
func (e *Encoder) SetReplacer(r Replacer) {
	e.replacer = r
}
