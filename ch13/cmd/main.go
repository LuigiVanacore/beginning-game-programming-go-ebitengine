package main

import (
	"log"

	"book/code/ch13"
)

func main() {
	if err := game.NewApp().Run(); err != nil {
		log.Fatal(err)
	}
}
