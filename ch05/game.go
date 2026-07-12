package game


import (
	. "book/code/ch05/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	engine *Engine
	player *Player
}

func NewGame() *Game {
	engine := NewEngine()
	rm := engine.ResourceManager()
	loadTextures(rm)
	world := engine.World()
	registerPlayerInput(engine.Input())
	setupTilemap(world, rm, defaultSceneConfig)
	player := NewPlayer(world, rm)
	world.Camera().SetFollow(player)
	return &Game{engine: engine, player: player}
}

func (g *Game) Update() error {
	g.player.UpdateMovement(g.engine.Input())
	return g.engine.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.engine.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.engine.Layout(outsideWidth, outsideHeight)
}
