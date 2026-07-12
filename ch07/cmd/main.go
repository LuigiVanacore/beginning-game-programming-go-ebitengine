package main

import (
	"log"

	"book/code/ch07"
)

func main() {
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
