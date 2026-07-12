package assets

// Paths are relative to the embedded assets/ directory (see FS in embed.go).
// Every runtime texture now lives under sprites/.
const (
	// Floor tileset: spritesheet.png stacks three 16x16 floor variants vertically
	// (1px spacing) — plain, crates, and stones. FloorTileIndex maps floor.map
	// indices 0/1/2 to cells (0,0)/(0,1)/(0,2). Same file and key as Chapter 10.
	Spritesheet = "sprites/spritesheet.png"

	Player      = "sprites/player.png"
	Enemy       = "sprites/enemy.png"
	Cursor      = "sprites/cursor.png"
	BloodyKnife = "sprites/bloody_knife.png"
	Potion      = "sprites/potion.png"
	SacredBook  = "sprites/sacred_book.png"
	FlyingAxe   = "sprites/flying_axe.png"

	// Upgrade icons and bonus items.
	Armor  = "sprites/defense_armor.png"
	Boots  = "sprites/speed_boots.png"
	Gem    = "sprites/experience_gem.png"
	Skull  = "sprites/cursed_skull.png"
	Ring   = "sprites/power_ring.png"
	Shield = "sprites/holy_shield.png"

	// Monster sprites, ordered from the weakest to the strongest kind.
	Ghost      = "sprites/ghost.png"
	Spider     = "sprites/spider.png"
	Bat        = "sprites/bat.png"
	DarkWizard = "sprites/dark_wizard.png"
	Cyclops    = "sprites/cyclops.png"

	FloorMap = "floor.map"

	// Sound effects and music under audio/ (new in ch12). Each is a small 16-bit
	// mono WAV; the audio system decodes and registers them under SoundID keys.
	HitSound       = "audio/hit.wav"
	PickupSound    = "audio/pickup.wav"
	LevelUpSound   = "audio/levelup.wav"
	ClickSound     = "audio/click.wav"
	MusicTrack     = "audio/music.wav"
	EnemyGrowSound = "audio/enemy_grow.wav"
)
