package jsonreplace_test

import (
	"encoding/json"
	"testing"

	"github.com/mashiike/jsonreplace"
	"github.com/stretchr/testify/require"
)

func TestEncoderErrAbortReplacer(t *testing.T) {
	org := Organization{
		Leader: Person{
			Email: "admin@example.com",
			Name:  "Tarou Yamada",
			Age:   30,
		},
		Members: []Person{
			{
				Email: "member1@example.com",
				Name:  "Hanako Tanaka",
				Age:   20,
			},
			{
				Email: "member2@exampl.com",
				Name:  "Jhon Smith",
				Age:   25,
			},
		},
	}
	mux := jsonreplace.NewReplaceMux()
	mux.ReplaceFunc(`{"type":"object","properties":{"age":{"type":"integer"}},"required":["age"]}`, func(raw json.RawMessage) (json.RawMessage, error) {
		var v map[string]interface{}
		if err := json.Unmarshal(raw, &v); err != nil {
			return nil, err
		}
		if num, ok := v["age"].(float64); ok {
			num -= 5
			if num < 20 {
				num = 20
			}
			v["age"] = num
		}
		bs, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return bs, jsonreplace.ErrAbortReplacer
	})
	mux.ReplaceFunc(`{"type":"string","format":"email"}`, func(raw json.RawMessage) (json.RawMessage, error) {
		return json.RawMessage(`"***********@example.com"`), nil
	})
	bs, err := jsonreplace.MarshalIndent(org, mux, "", "  ")
	require.NoError(t, err)
	t.Log(string(bs))
	require.JSONEq(t, `{
		"leader": {
		  "age": 25,
		  "email": "admin@example.com",
		  "name": "Tarou Yamada"
		},
		"members": [
		  {
			"age": 20,
			"email": "member1@example.com",
			"name": "Hanako Tanaka"
		  },
		  {
			"age": 20,
			"email": "member2@exampl.com",
			"name": "Jhon Smith"
		  }
		]
	  }`, string(bs))
}
