# Ninesweeper

A hybrid of Sudoku and Minesweeper built with [Ebiten](https://ebiten.org/) in Go.

## The Concept

A 9x9 Sudoku grid where some cells are mines. Use **both** Minesweeper logic (adjacency counts) and Sudoku logic (row/column/box constraints) to figure out which cells are safe and what digit goes where.

For the math behind why this works, see [HOW_IT_WORKS.md](HOW_IT_WORKS.md).

## Features

- **Title screen** with difficulty selector and best times
- **4 difficulty levels**: Easy, Medium, Hard, Expert
- **High score tracking** — best times saved per difficulty
- **Scores screen** — view your top 5 times per difficulty
- **WASM build** — play in the browser, deploy anywhere

## How to Play

**The Goal:** Reveal all non-mine cells and fill them with the correct Sudoku digit.

**Controls:**

| Input | Action |
|-------|--------|
| Left Click (hidden cell) | Reveal it |
| Right Click (hidden cell) | Toggle mine flag |
| Left Click (revealed cell) | Select for digit entry |
| 1-9 keys / numpad buttons | Place a digit |
| Backspace / Delete | Clear digit |
| R | Restart (same difficulty) |
| Esc | Back to menu |

**Visual Guide:**
- **Green cells** / green digits = Pre-revealed givens (safe, locked)
- **White cells** = Revealed safe cells (enter your guess)
- **Yellow highlight** = Selected cell
- **Blue/gray cells** = Hidden (click to reveal or flag)
- **Orange cells** = Flagged as mine
- **Small top-right number** = Adjacent mine count
- **Large center number** = Sudoku digit (green=given, blue=guess, red=conflict)

## Difficulty Levels

| Level | Mines | Givens | Vibe |
|-------|-------|--------|------|
| Easy | 8 | 25 | Chill |
| Medium | 12 | 20 | Fair fight |
| Hard | 15 | 15 | Sweaty |
| Expert | 18 | 12 | Pain |

## Build & Run (Desktop)

```powershell
# Install Go (1.22+) from https://go.dev/dl/

git clone https://github.com/saedarm/ninesweeper.git
cd ninesweeper
go mod tidy
go run .
```

## Build & Run (WASM / Browser)

### Step 1: Build the WASM binary

```powershell
.\build.ps1
```

This does three things:
1. Compiles the game to `build/ninesweeper.wasm` (using `GOOS=js GOARCH=wasm`)
2. Copies Go's `wasm_exec.js` bridge file into `build/`
3. Generates `build/index.html` with a dark-themed wrapper page

### Step 2: Test locally

```powershell
cd build
python -m http.server 8080
```

Open [http://localhost:8080](http://localhost:8080) in your browser. You should see the Ninesweeper title screen. If it works locally, it'll work deployed.

## Deploy to Azure (Static Web Apps)

The `build/` folder is just three static files. No server, no container, no Docker. Azure Static Web Apps is the easiest way to host this.

### Prerequisites

- [Azure CLI](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli-windows) installed
- An Azure account ([free tier](https://azure.microsoft.com/en-us/free/) works fine)
- The WASM build completed (Step 1 above)

### Step 1: Log in to Azure

```powershell
az login
```

This opens a browser window for authentication.

### Step 2: Create a resource group (skip if you already have one)

```powershell
az group create --name ninesweeper-rg --location eastus
```

### Step 3: Create the Static Web App

The easiest path is to deploy from a GitHub repo. Push your code (including the `build/` folder) to `github.com/saedarm/ninesweeper`, then:

```powershell
az staticwebapp create `
  --name ninesweeper `
  --resource-group ninesweeper-rg `
  --source https://github.com/saedarm/ninesweeper `
  --location eastus `
  --branch main `
  --app-location "/build" `
  --output-location "" `
  --login-with-github
```

This will:
1. Prompt you to authenticate with GitHub
2. Create a GitHub Actions workflow in your repo
3. Deploy the contents of `/build` as your static site

### Step 4: Get your URL

```powershell
az staticwebapp show --name ninesweeper --resource-group ninesweeper-rg --query "defaultHostname" -o tsv
```

That's your live URL. Done.

### Alternative: Azure Portal (no CLI)

1. Go to [portal.azure.com](https://portal.azure.com) → Create a resource → **Static Web App**
2. Connect your GitHub repo (`saedarm/ninesweeper`)
3. Under **Build Details**:
   - **Build Preset**: Custom
   - **App location**: `/build`
   - **Output location**: (leave blank)
4. Click **Review + Create** → **Create**
5. Azure auto-generates a GitHub Actions workflow and deploys

The key thing: set **Build Preset to Custom** and point **App location to `/build`**. Azure doesn't need to build anything — the WASM is already compiled. It just serves the files.

## Project Structure

```
ninesweeper/
├── main.go              # Game loop, screen states, input handling
├── board/
│   └── board.go         # Core game state: cells, mines, reveal, win/loss
├── sudoku/
│   └── sudoku.go        # Sudoku puzzle generator (backtracking)
├── render/
│   └── render.go        # All drawing: menu, game, scores, grid, UI
├── config/
│   └── config.go        # Difficulty level definitions
├── scores/
│   └── scores.go        # High score tracking (JSON file persistence)
├── build.ps1            # WASM build script (PowerShell)
├── HOW_IT_WORKS.md      # Math explainer
├── go.mod
└── README.md
```

## License

MIT
