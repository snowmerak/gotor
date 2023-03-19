package main

import (
	"fmt"
	"github.com/snowmerak/gotor/actor"
)

func main() {
	rs := actor.Generate("actor", "Actor", map[string]string{"a": "int", "b": "string"})
	fmt.Println(string(rs))
}
