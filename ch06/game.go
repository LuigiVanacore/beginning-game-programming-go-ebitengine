package game


import (
	. "book/code/ch06/internal/core"
	en "book/code/ch06/enemy"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game holds the session state: world, nodes, HUD.
type Game struct {
	engine   *Engine
	player   *Player
	enemy    *en.Enemy
	gameOver bool
	hud      *HUD
}

func NewGame() *Game {
	engine := NewEngine()
	rm := engine.ResourceManager()
	loadTextures(rm)

	world := engine.World()
	registerPlayerInput(engine.Input())
	setupTilemap(world, rm, defaultSceneConfig)

	player := NewPlayer(engine)
	enemy := en.NewEnemyFarFromPlayer(engine, 0, 0)
	world.Camera().SetFollow(player)

	game := &Game{
		engine:   engine,
		player:   player,
		enemy:    enemy,
		gameOver: false,
		hud:      NewHUD(),
	}
	wirePlayerGameOver(player, game)
	return game
}

func (g *Game) Update() error {
	if g.gameOver {
		return g.engine.Update()
	}
	g.player.Update()
	playerPos := g.player.GetWorldPosition()
	g.enemy.Update(playerPos.X(), playerPos.Y())
	g.engine.CollisionManager().CheckCollision()
	return g.engine.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.engine.Draw(screen)
	if g.gameOver {
		g.hud.DrawGameOver(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.engine.Layout(outsideWidth, outsideHeight)
}
