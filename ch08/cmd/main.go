package main

import (
	"log"

	"book/code/ch08"
)

func main() {
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
