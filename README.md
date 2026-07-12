# Beginning Game Programming with Go and Ebitengine — Codice sorgente

Questo repository contiene il codice sorgente completo del libro **_Beginning Game
Programming with Go and Ebitengine_**.

Il libro insegna, passo dopo passo, come programmare un videogioco 2D in stile
**"bullet-heaven"** (il genere reso popolare da *Vampire Survivors*) usando il
linguaggio **Go** e il framework **[Ebitengine](https://ebitengine.org/)**.

Il gioco che si costruisce lungo i capitoli si chiama **Gopher Survivor**: il
giocatore controlla un Gopher che si muove su una mappa infinita, viene assalito
da ondate di nemici sempre più numerose, spara e attiva armi in automatico,
raccoglie esperienza, sale di livello e sceglie potenziamenti, fino allo scontro
finale con il boss.

---

## A chi è rivolto

- Sviluppatori che conoscono le basi di Go e vogliono imparare il game programming.
- Chi vuole capire come si costruisce **da zero** un piccolo motore 2D
  (scene graph, gestione risorse, input, collisioni, camera, UI, audio, stati di gioco).
- Chi è curioso del genere "bullet-heaven" e vuole vederne l'implementazione reale.

Non serve alcuna esperienza pregressa con Ebitengine.

---

## Requisiti

| Strumento | Versione |
|-----------|----------|
| **Go** | 1.22.0 o superiore |
| **Ebitengine** | v2.8.6 |

Dipendenze principali (gestite via Go Modules, vedi i vari `go.mod`):

- `github.com/hajimehoshi/ebiten/v2` — il framework di gioco
- `golang.org/x/image` — font e utility grafiche (dai capitoli con UI/testo in poi)

> **Nota per Linux:** Ebitengine richiede alcune librerie di sistema (OpenGL, ALSA,
> X11/Wayland). Consulta la [pagina di installazione ufficiale](https://ebitengine.org/en/documents/install.html)
> se la compilazione fallisce per dipendenze native mancanti.

---

## Struttura del repository

Ogni capitolo del libro è una **cartella autonoma e autoconsistente** (`ch01` … `ch13`),
con il proprio `go.mod`. Questo permette di compilare ed eseguire ogni capitolo in
modo indipendente e di confrontare facilmente lo stato del codice tra un passo e
il successivo.

```
go_game_dev_book_code/
├── ch01/               # un capitolo = un modulo Go a sé stante
│   ├── go.mod          #   modulo "book/code/ch01"
│   ├── cmd/main.go     #   entry point (func main)
│   ├── assets/         #   sprite, tileset, audio, font (embedded via go:embed)
│   ├── internal/core/  #   il "mini-motore" costruito nel libro
│   └── *.go            #   codice specifico del capitolo (player, nemici, armi, ...)
├── ch02/
├── ...
└── ch13/
```

### Convenzioni ricorrenti

Dai primi capitoli in poi si consolida una struttura comune, riusata e ampliata
capitolo per capitolo:

- **`cmd/main.go`** — punto di ingresso; costruisce l'app/gioco e chiama `ebiten.RunGame`.
- **`run.go` / `app.go`** — configurazione della finestra e avvio del game loop.
- **`internal/core/`** — il motore riutilizzabile: scene graph, nodi, trasformazioni,
  gestione risorse, input, collisioni, camera, tilemap, timer, audio, macchina a stati.
- **`assets/`** — risorse incorporate nel binario tramite `//go:embed` (nessun file
  esterno da distribuire con l'eseguibile).
- **`enemy/`, `pickups/`, `ui/`** — pacchetti che raggruppano le entità di gioco man
  mano che vengono introdotte.

Il cuore dell'architettura è un **motore basato su scene graph** (ispirato ai game
engine "a nodi"): un `Engine` possiede un `World`; il `World` contiene un albero di
`Node`/`Node2D` che vengono aggiornati e disegnati ad ogni frame, secondo il classico
ciclo `Update` → `Draw` → `Layout` di Ebitengine.

---

## I capitoli

Ogni capitolo aggiunge un tassello sopra il precedente. Di seguito il percorso
completo, con il titolo mostrato nella finestra di gioco:

| Cap. | Titolo | Cosa si impara |
|------|--------|----------------|
| **01** | Hello Ebiten - Go Gopher | Primo programma Ebitengine: finestra, game loop, disegno di uno sprite (il Gopher). |
| **02** | Scene Graph Framework | Fondamenta del motore: `Engine`, `World`, `Node`/`Node2D`, trasformazioni, vettori 2D, scene graph. |
| **03** | Resource Manager, Layers, Sprites | Gestione centralizzata delle risorse, ordinamento in layer, sprite. |
| **04** | Input and Player Movement | Sistema di input (action map), movimento del giocatore. |
| **05** | Tileset, Tilemap, Camera | Tileset, tilemap infinita e una camera che segue il giocatore. |
| **06** | Enemy and Collisions | Introduzione dei nemici e sistema di collisioni (collider, maschere, manager). |
| **07** | Weapons and Projectiles | Armi automatiche, proiettili e **object pooling** per riutilizzare le entità. |
| **08** | UI, Health, XP, Level Up | Interfaccia: barra della salute, barra XP, HUD, popup, schermata di *game over*, salita di livello. |
| **09** | Potions, Sacred Book, Holy Shield | Nuovi pickup (pozioni) e armi orbitanti/ad area (Sacred Book, Holy Shield). |
| **10** | Weapon Upgrade UI | Nuove armi (Flying Axe) e schermata di scelta dei potenziamenti. |
| **11** | Gopher Survivor | Consolidamento del gameplay: sistema di upgrade (bonus, armi), utility, loop di gioco completo. |
| **12** | Gopher Survivor — State machine | **Macchina a stati** dell'applicazione (menu, opzioni, gioco, pausa), audio, difficoltà e tipi di nemici. |
| **13** | Gopher Survivor (finale) | Rifiniture finali: particellari, testo fluttuante (danni), tutte le funzionalità integrate. |

> Il capitolo **13** rappresenta la versione finale e completa del gioco.

---

## Come compilare ed eseguire

Ogni capitolo si esegue in modo indipendente. Dalla cartella del capitolo desiderato:

```bash
# esempio: eseguire il capitolo 1
cd ch01
go run ./cmd

# esempio: eseguire il gioco completo (capitolo finale)
cd ch13
go run ./cmd
```

Per produrre un eseguibile:

```bash
cd ch13
go build -o gopher-survivor ./cmd
```

Poiché gli asset sono **incorporati** nel binario tramite `//go:embed`, l'eseguibile
risultante è autonomo e può essere distribuito senza cartelle di risorse aggiuntive.

### Comandi di gioco (capitolo finale)

- **Movimento:** WASD / frecce direzionali
- Le armi si attivano **automaticamente**; l'obiettivo è sopravvivere il più a lungo
  possibile e potenziare il proprio Gopher salendo di livello.

---

## Come studiare il codice con il libro

1. Leggi il capitolo del libro.
2. Apri la cartella `chNN` corrispondente ed eseguila.
3. Confronta il codice con quello del capitolo precedente per vedere esattamente
   **cosa è cambiato** (l'organizzazione a moduli separati rende il diff immediato).

Questo approccio "una cartella per capitolo" è pensato apposta come materiale
didattico: puoi tornare a qualunque tappa del percorso senza dover ripristinare
manualmente lo stato del progetto.

---

## Licenza

Il codice sorgente di questo repository è distribuito sotto **Licenza MIT** — vedi
il file [LICENSE](LICENSE).

Il **testo del libro** _Beginning Game Programming with Go and Ebitengine_ e le sue
illustrazioni **non** sono coperti dalla licenza MIT e restano protetti dai rispettivi
diritti d'autore.

Gli asset di terze parti (ad es. l'immagine del Go Gopher, font, eventuali risorse
audio/grafiche) restano soggetti alle rispettive licenze originali. Il logo/mascotte
**Go Gopher** è stato creato da Renée French ed è distribuito sotto licenza
[Creative Commons Attribution 4.0](https://creativecommons.org/licenses/by/4.0/).

---

## Crediti

- **Autore:** Luigi Vanacore
- **Framework:** [Ebitengine](https://ebitengine.org/) di Hajime Hoshi
- **Linguaggio:** [Go](https://go.dev/)
