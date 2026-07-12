package main

import (
	"log"

	"book/code/ch02"
)

func main() {
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
