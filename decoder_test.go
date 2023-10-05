package jsonreplace_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mashiike/jsonreplace"
	"github.com/stretchr/testify/require"
)

type DecoderTestCase struct {
	CaseName      string
	Input         string
	Expected      interface{}
	Error         string
	DecoderConfig func(*testing.T, *jsonreplace.Decoder)
}

func (dtc *DecoderTestCase) Run(t *testing.T) {
	decoder := jsonreplace.NewDecoder(strings.NewReader(dtc.Input))
	if dtc.DecoderConfig != nil {
		dtc.DecoderConfig(t, decoder)
	}
	var v interface{}
	err := decoder.Decode(&v)
	if dtc.Error != "" {
		require.EqualError(t, err, dtc.Error)
		return
	}
	require.NoError(t, err)
	require.EqualValues(t, dtc.Expected, v)
}

func TestDecoder__NoReplace(t *testing.T) {
	cases := []DecoderTestCase{
		{
			CaseName: "null",
			Input:    "null",
			Expected: nil,
		},
		{
			CaseName: "empty object",
			Input:    "{}",
			Expected: map[string]interface{}{},
		},
		{
			CaseName: "empty array",
			Input:    "[]",
			Expected: []interface{}{},
		},
		{
			CaseName: "empty string",
			Input:    `""`,
			Expected: "",
		},
		{
			CaseName: "object with email",
			Input:    `{"email":"baz@example.com","foo":"bar"}`,
			Expected: map[string]interface{}{
				"foo":   "bar",
				"email": "baz@example.com",
			},
		},
		{
			CaseName: "array with email",
			Input:    `["foo","bar","baz@example.com"]`,
			Expected: []interface{}{"foo", "bar", "baz@example.com"},
		},
	}
	for _, c := range cases {
		t.Run(c.CaseName, c.Run)
	}
}

func TestDecoder__Replace(t *testing.T) {
	decoderConfig := func(t *testing.T, decoder *jsonreplace.Decoder) {
		mux := jsonreplace.NewReplaceMux()
		mux.ReplaceFunc(`{"type":"string", "format": "email"}`, func(bs json.RawMessage) (json.RawMessage, error) {
			return json.RawMessage(`"********@example.com"`), nil
		})
		decoder.SetReplacer(mux)
	}

	cases := []DecoderTestCase{
		{
			CaseName:      "null",
			Input:         "null",
			Expected:      nil,
			DecoderConfig: decoderConfig,
		},
		{
			CaseName:      "empty object",
			Input:         "{}",
			Expected:      map[string]interface{}{},
			DecoderConfig: decoderConfig,
		},
		{
			CaseName:      "empty array",
			Input:         "[]",
			Expected:      []interface{}{},
			DecoderConfig: decoderConfig,
		},
		{
			CaseName:      "empty string",
			Input:         `""`,
			Expected:      "",
			DecoderConfig: decoderConfig,
		},
		{
			CaseName: "object with email",
			Input:    `{"email":"baz@example.com","foo":"bar"}`,
			Expected: map[string]interface{}{
				"foo":   "bar",
				"email": "********@example.com",
			},
			DecoderConfig: decoderConfig,
		},
		{
			CaseName:      "array with email",
			Input:         `["foo","bar","baz@example.com"]`,
			Expected:      []interface{}{"foo", "bar", "********@example.com"},
			DecoderConfig: decoderConfig,
		},
	}
	for _, c := range cases {
		t.Run(c.CaseName, c.Run)
	}
}
