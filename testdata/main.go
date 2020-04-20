package main

import (
	"encoding/json"
	"fmt"
	"log"
)

//go:generate go run github.com/darkclainer/typedjson -interface Data *Foo *Bar
type Data interface {
	typedjson(*DataTyped) string
}

type Foo struct {
	Name  string
	Value int
}

type Bar struct {
	Name1   string
	Another []string
}

func main() {
	one := Foo{
		Name: "first",
	}

	encoded, err := json.Marshal(DataTyped{Data: &one})
	if err != nil {
		log.Fatal("Marshal json: ", err)
	}

	fmt.Println(string(encoded))

	var decoded DataTyped

	err = json.Unmarshal(encoded, &decoded)
	if err != nil {
		log.Fatal("Unmarshal json: ", err)
	}

	oneDecoded, ok := decoded.Data.(*Foo)
	if !ok {
		log.Fatal("type not match")
	}

	fmt.Printf("%#v\n", oneDecoded)
}
