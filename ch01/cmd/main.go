package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"

	"book/code/ch01"
	"book/code/ch01/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	img, _, err := image.Decode(bytes.NewReader(assets.GopherPNG))
	if err != nil {
		log.Fatalf("failed to decode embedded gopher image: %v", err)
	}
	eimg := ebiten.NewImageFromImage(img)

	g := game.NewGame(eimg)

	ebiten.SetWindowSize(game.Settings.ScreenWidth, game.Settings.ScreenHeight)
	ebiten.SetWindowTitle("Chapter 1: Hello Ebiten - Go Gopher")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
