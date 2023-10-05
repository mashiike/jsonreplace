# jsonreplace
This package provides a JSON utility function for replacing values in JSON.

[![GoDoc](https://godoc.org/github.com/mashiike/jsonreplace?status.svg)](https://godoc.org/github.com/mashiike/jsonreplace)
[![Go Report Card](https://goreportcard.com/badge/github.com/mashiike/jsonreplace)](https://goreportcard.com/report/github.com/mashiike/jsonreplace)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Overview

This package provides a JSON utility function for replacing values in JSON. 
It is enabled to replace JSON valuse on matched JSON Schema.  

```go 
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mashiike/jsonreplace"
)

type Person struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type Organization struct {
	Leader  Person   `json:"leader"`
	Members []Person `json:"members"`
}

func main() {
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
	jsonreplace.ReplaceFunc(`{"type":"object","properties":{"age":{"type":"integer"}},"required":["age"]}`, func(raw json.RawMessage) (json.RawMessage, error) {
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
		return json.Marshal(v)
	})
	jsonreplace.ReplaceFunc(`{"type":"string","format":"email"}`, func(raw json.RawMessage) (json.RawMessage, error) {
		return json.RawMessage(`"***********@example.com"`), nil
	})
	bs, err := jsonreplace.MarshalIndent(org, nil, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
    fmt.Println(string(bs))
	// Output:
	// {
    //   "leader": {
    //     "email": "***********@example.com",
    //     "name": "Tarou Yamada",
    //     "age": 25
    //   },
    //   "members": [
    //     {
    //       "email": "***********@example.com",
    //       "name": "Hanako Tanaka",
    //       "age": 20
    //     },
    //     {
    //       "email": "***********@example.com",
    //       "name": "Jhon Smith",
    //       "age": 20
    //     }
    //   ]
    // }
}

```

this example replace email address and age value in JSON.
email address masked by `***********@example.com` and age value is decreased by 5 and minimum value is 20.

## Installation

```bash
go get github.com/mashiike/slogutils
```

## License
This project is licensed under the MIT License - see the LICENSE(./LICENCE) file for details.

## Contribution
Contributions, bug reports, and feature requests are welcome. Pull requests are also highly appreciated. For more details, please
