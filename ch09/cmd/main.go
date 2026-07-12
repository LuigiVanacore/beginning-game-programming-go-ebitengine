package main

import (
	"log"

	"book/code/ch09"
)

func main() {
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
