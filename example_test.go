package jsonreplace_test

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

func ExampleMarshal() {
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
	var masked Organization
	if err := json.Unmarshal(bs, &masked); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Leader: Name=%s, Email=%s, Age=%d\n", masked.Leader.Name, masked.Leader.Email, masked.Leader.Age)
	for i, member := range masked.Members {
		fmt.Printf("Member%d: Name=%s, Email=%s, Age=%d\n", i+1, member.Name, member.Email, member.Age)
	}
	fmt.Println()
	// Output:
	//Leader: Name=Tarou Yamada, Email=***********@example.com, Age=25
	//Member1: Name=Hanako Tanaka, Email=***********@example.com, Age=20
	//Member2: Name=Jhon Smith, Email=***********@example.com, Age=20
}
