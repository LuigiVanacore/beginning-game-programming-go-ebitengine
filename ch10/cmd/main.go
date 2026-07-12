package main

import (
	"log"

	"book/code/ch10"
)

func main() {
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
