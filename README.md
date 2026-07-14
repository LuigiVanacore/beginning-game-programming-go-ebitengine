# Beginning Game Programming with Go — Source Code

This repository contains the complete source code for the book
[**_Beginning Game Programming with Go_**](https://leanpub.com/gameprogramminggolang).

The book teaches, step by step, how to program a 2D **"bullet-heaven"** game (the
genre popularized by *Vampire Survivors*) using the **Go** language and the
**[Ebitengine](https://ebitengine.org/)** framework.

The game you build throughout the chapters is called **Gopher Survivor**: you
control a Gopher moving across an infinite map, get swarmed by ever-growing waves
of enemies, fire and trigger weapons automatically, collect experience, level up
and pick upgrades, all the way to the final boss encounter.

---

## Who this is for

- Developers who know the basics of Go and want to learn game programming.
- Anyone who wants to understand how to build a small 2D engine **from scratch**
  (scene graph, resource management, input, collisions, camera, UI, audio, game states).
- Anyone curious about the "bullet-heaven" genre who wants to see a real implementation.

No prior experience with Ebitengine is required.

---

## Requirements

| Tool | Version |
|------|---------|
| **Go** | 1.22.0 or later |
| **Ebitengine** | v2.8.6 |

Main dependencies (managed via Go Modules, see the individual `go.mod` files):

- `github.com/hajimehoshi/ebiten/v2` the game framework
- `golang.org/x/image` fonts and graphics utilities (from the chapters with UI/text onward)

> **Note for Linux:** Ebitengine requires some system libraries (OpenGL, ALSA,
> X11/Wayland). See the [official installation page](https://ebitengine.org/en/documents/install.html)
> if the build fails because of missing native dependencies.

---

## Repository structure

Each chapter of the book is a **self-contained, standalone folder** (`ch01` … `ch13`),
with its own `go.mod`. This lets you build and run each chapter independently and
easily compare the state of the code from one step to the next.

```
go_game_dev_book_code/
├── ch01/               # one chapter = one standalone Go module
│   ├── go.mod          #   module "book/code/ch01"
│   ├── cmd/main.go     #   entry point (func main)
│   ├── assets/         #   sprites, tilesets, audio, fonts (embedded via go:embed)
│   ├── internal/core/  #   the "mini-engine" built throughout the book
│   └── *.go            #   chapter-specific code (player, enemies, weapons, ...)
├── ch02/
├── ...
└── ch13/
```

### Recurring conventions

From the early chapters onward, a common structure takes shape and is reused and
expanded chapter by chapter:

- **`cmd/main.go`** — entry point; builds the app/game and calls `ebiten.RunGame`.
- **`run.go` / `app.go`** — window configuration and game-loop startup.
- **`internal/core/`** — the reusable engine: scene graph, nodes, transforms,
  resource management, input, collisions, camera, tilemap, timers, audio, state machine.
- **`assets/`** — resources embedded into the binary via `//go:embed` (no external
  files to ship alongside the executable).
- **`enemy/`, `pickups/`, `ui/`** — packages that group the game entities as they
  are introduced.

The core of the architecture is a **scene-graph-based engine** (inspired by
node-based game engines): an `Engine` owns a `World`; the `World` holds a tree of
`Node`/`Node2D` objects that are updated and drawn every frame, following
Ebitengine's classic `Update` → `Draw` → `Layout` cycle.

---

## The chapters

Each chapter builds one more piece on top of the previous one. Below is the full
path, with the title shown in the game window:

| Ch. | Title | What you learn |
|-----|-------|----------------|
| **01** | Hello Ebiten - Go Gopher | Your first Ebitengine program: window, game loop, drawing a sprite (the Gopher). |
| **02** | Scene Graph Framework | Engine foundations: `Engine`, `World`, `Node`/`Node2D`, transforms, 2D vectors, scene graph. |
| **03** | Resource Manager, Layers, Sprites | Centralized resource management, layer ordering, sprites. |
| **04** | Input and Player Movement | Input system (action map), player movement. |
| **05** | Tileset, Tilemap, Camera | Tileset, infinite tilemap and a camera that follows the player. |
| **06** | Enemy and Collisions | Introducing enemies and the collision system (colliders, masks, manager). |
| **07** | Weapons and Projectiles | Automatic weapons, projectiles and **object pooling** to reuse entities. |
| **08** | UI, Health, XP, Level Up | Interface: health bar, XP bar, HUD, popups, *game over* screen, leveling up. |
| **09** | Potions, Sacred Book, Holy Shield | New pickups (potions) and orbiting/area weapons (Sacred Book, Holy Shield). |
| **10** | Weapon Upgrade UI | New weapons (Flying Axe) and the upgrade-selection screen. |
| **11** | Gopher Survivor | Gameplay consolidation: upgrade system (bonuses, weapons), utilities, complete game loop. |
| **12** | Gopher Survivor — State machine | Application **state machine** (menu, options, game, pause), audio, difficulty and enemy types. |
| **13** | Gopher Survivor (final) | Final polish: particles, floating (damage) text, all features integrated. |

> Chapter **13** is the final, complete version of the game.

---

## Building and running

Each chapter runs independently. From the folder of the chapter you want:

```bash
# example: run chapter 1
cd ch01
go run ./cmd

# example: run the complete game (final chapter)
cd ch13
go run ./cmd
```

To produce an executable:

```bash
cd ch13
go build -o gopher-survivor ./cmd
```

Because the assets are **embedded** into the binary via `//go:embed`, the resulting
executable is self-contained and can be distributed without any extra resource folders.

### Controls (final chapter)

- **Movement:** WASD / arrow keys
- Weapons fire **automatically**; the goal is to survive as long as possible and
  power up your Gopher by leveling up.

---

## How to study the code with the book

1. Read the book chapter.
2. Open the matching `chNN` folder and run it.
3. Compare the code with the previous chapter to see exactly **what changed** (the
   one-module-per-chapter layout makes the diff immediate).

This "one folder per chapter" approach is designed specifically as teaching
material: you can go back to any point in the journey without having to manually
restore the project state.

---

## License

The source code in this repository is released under the **MIT License** — see the
[LICENSE](LICENSE) file.

The **text** of the book _Beginning Game Programming with Go and Ebitengine_ and its
illustrations are **not** covered by the MIT License and remain protected by their
respective copyrights.

Third-party assets (e.g. the Go Gopher image, fonts, any audio/graphic resources)
remain subject to their own original licenses. The **Go Gopher** logo/mascot was
created by Renée French and is distributed under the
[Creative Commons Attribution 4.0](https://creativecommons.org/licenses/by/4.0/) license.

---

## Credits

- **Author:** Luigi Vanacore
- **Framework:** [Ebitengine](https://ebitengine.org/) by Hajime Hoshi
- **Language:** [Go](https://go.dev/)
