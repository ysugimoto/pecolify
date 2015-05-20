package main

import (
	"fmt"
	"github.com/ysugimoto/pecolify"
)

func main() {
	// Instanciate pecolify
	pf := pecolify.New()

	// Pass data to pecolify
	data := []string{
		"foo",
		"bar",
		"baz",
	}

	// pecolify!
	selected, err := pf.Transform(data)
	if err != nil {
		fmt.Printf("Error was occured: %v\n", err)
		return
	}

	fmt.Printf("Selected from peco: %s\n", selected)
}
