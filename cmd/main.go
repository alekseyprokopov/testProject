package main

import (
	"fmt"
	"log"
	"testProject/app/apiserver"
)

func main() {
	s, err := apiserver.New()
	if err != nil {
		log.Fatal(fmt.Errorf("can't create apiserver: %w", err))
	}

	s.Start()
}
