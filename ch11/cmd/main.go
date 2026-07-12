package main

import (
	"log"

	"book/code/ch11"
)

func main() {
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
