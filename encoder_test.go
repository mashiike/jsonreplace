package jsonreplace_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/mashiike/jsonreplace"
	"github.com/stretchr/testify/require"
)

type EncoderTestCase struct {
	CaseName      string
	Value         interface{}
	Expected      string
	Error         string
	EncoderConfig func(*testing.T, *jsonreplace.Encoder)
}

func (etc *EncoderTestCase) Run(t *testing.T) {
	var buf bytes.Buffer
	encoder := jsonreplace.NewEncoder(&buf)
	if etc.EncoderConfig != nil {
		etc.EncoderConfig(t, encoder)
	}
	err := encoder.Encode(etc.Value)
	if etc.Error != "" {
		require.EqualError(t, err, etc.Error)
		return
	}
	require.NoError(t, err)
	require.JSONEq(t, string(etc.Expected), buf.String())
}

func TestEncoder__NoReplace(t *testing.T) {
	cases := []EncoderTestCase{
		{
			CaseName: "nil",
			Value:    nil,
			Expected: "null",
		},
		{
			CaseName: "empty object",
			Value:    map[string]interface{}{},
			Expected: "{}",
		},
		{
			CaseName: "empty array",
			Value:    []interface{}{},
			Expected: "[]",
		},
		{
			CaseName: "empty string",
			Value:    "",
			Expected: `""`,
		},
		{
			CaseName: "object with email",
			Value: map[string]interface{}{
				"foo":   "bar",
				"email": "baz@example.com",
			},
			Expected: `{"foo":"bar","email":"baz@example.com"}`,
		},
		{
			CaseName: "array with email",
			Value:    []interface{}{"foo", "bar", "baz@example.com"},
			Expected: `["foo","bar","baz@example.com"]`,
		},
	}
	for _, etc := range cases {
		t.Run(etc.CaseName, etc.Run)
	}
}

func TestEncoder__WithReplace(t *testing.T) {
	encoderConfig := func(t *testing.T, encoder *jsonreplace.Encoder) {
		mux := jsonreplace.NewReplaceMux()
		mux.ReplaceFunc(`{"type":"string", "format": "email"}`, func(bs json.RawMessage) (json.RawMessage, error) {
			return json.RawMessage(`"********@example.com"`), nil
		})
		encoder.SetReplacer(mux)
	}
	cases := []EncoderTestCase{
		{
			CaseName:      "nil",
			Value:         nil,
			Expected:      "null",
			EncoderConfig: encoderConfig,
		},
		{
			CaseName:      "empty object",
			Value:         map[string]interface{}{},
			Expected:      "{}",
			EncoderConfig: encoderConfig,
		},
		{
			CaseName:      "empty array",
			Value:         []interface{}{},
			Expected:      "[]",
			EncoderConfig: encoderConfig,
		},
		{
			CaseName:      "empty string",
			Value:         "",
			Expected:      `""`,
			EncoderConfig: encoderConfig,
		},
		{
			CaseName:      "email string",
			Value:         "baz@example.com",
			Expected:      `"********@example.com"`,
			EncoderConfig: encoderConfig,
		},
		{
			CaseName: "object with email",
			Value: map[string]interface{}{
				"foo":   "bar",
				"email": "baz@example.com",
			},
			Expected:      `{"foo":"bar","email":"********@example.com"}`,
			EncoderConfig: encoderConfig,
		},
		{
			CaseName:      "array with email",
			Value:         []interface{}{"foo", "bar", "baz@example.com"},
			Expected:      `["foo","bar","********@example.com"]`,
			EncoderConfig: encoderConfig,
		},
	}
	for _, etc := range cases {
		t.Run(etc.CaseName, etc.Run)
	}
}
