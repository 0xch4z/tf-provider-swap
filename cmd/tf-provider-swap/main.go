package main

import (
	"log"

	swap "github.com/charliekenney23/tf-provider-swap"
)

func main() {
	if err := swap.Entrypoint(); err != nil {
		log.Fatal(err)
	}
}
