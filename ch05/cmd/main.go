package main

import (
	"log"

	"book/code/ch05"
)

func main() {
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
