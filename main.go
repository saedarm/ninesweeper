package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/saedarm/ninesweeper/board"
	"github.com/saedarm/ninesweeper/config"
	"github.com/saedarm/ninesweeper/render"
	"github.com/saedarm/ninesweeper/scores"
)

type Screen int

const (
	ScreenMenu Screen = iota
	ScreenPlaying
	ScreenScores
)

type Game struct {
	screen     Screen
	board      *board.Board
	rng        *rand.Rand
	difficulty config.Difficulty
	hoverR     int
	hoverC     int
	menuHover  int
	startTime  time.Time
	elapsed    int
	started    bool
	scored     bool // whether we've already recorded the score for this round
}

func NewGame() *Game {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Game{
		screen:    ScreenMenu,
		rng:       rng,
		hoverR:    -1,
		hoverC:    -1,
		menuHover: -1,
	}
}

func (g *Game) startGame(d config.Difficulty) {
	g.difficulty = d
	g.board = board.New(g.rng, d.Mines, d.Givens)
	g.started = false
	g.scored = false
	g.elapsed = 0
	g.startTime = time.Now()
	g.screen = ScreenPlaying
}

func (g *Game) restart() {
	g.startGame(g.difficulty)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return render.ScreenSize()
}

func (g *Game) Update() error {
	switch g.screen {
	case ScreenMenu:
		return g.updateMenu()
	case ScreenPlaying:
		return g.updatePlaying()
	case ScreenScores:
		return g.updateScores()
	}
	return nil
}

// --- MENU ---

func (g *Game) updateMenu() error {
	mx, my := ebiten.CursorPosition()
	g.menuHover = render.MenuButtonAt(mx, my)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		idx := render.MenuButtonAt(mx, my)
		if idx >= 0 && idx < len(config.All) {
			g.startGame(config.All[idx])
		}
	}

	// Keyboard shortcuts: 1-4 for difficulty
	for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4} {
		if inpututil.IsKeyJustPressed(k) && i < len(config.All) {
			g.startGame(config.All[i])
		}
	}

	// S for scores
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.screen = ScreenScores
	}

	return nil
}

// --- PLAYING ---

func (g *Game) updatePlaying() error {
	mx, my := ebiten.CursorPosition()
	g.hoverR, g.hoverC = render.CellAt(mx, my)

	// Timer
	if g.started && g.board.GameState == board.StatePlaying {
		g.elapsed = int(time.Since(g.startTime).Seconds())
		if g.elapsed > 999 {
			g.elapsed = 999
		}
	}

	// Record score on win (once)
	if g.board.GameState == board.StateWon && !g.scored {
		g.scored = true
		scores.Add(g.difficulty.Name, g.elapsed)
	}

	// Left click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.handleLeftClick(mx, my)
	}

	// Right click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.handleRightClick(mx, my)
	}

	// Keyboard digits
	for k := ebiten.Key1; k <= ebiten.Key9; k++ {
		if inpututil.IsKeyJustPressed(k) {
			g.placeDigitOnSelected(int(k-ebiten.Key1) + 1)
		}
	}
	for k := ebiten.KeyNumpad1; k <= ebiten.KeyNumpad9; k++ {
		if inpututil.IsKeyJustPressed(k) {
			g.placeDigitOnSelected(int(k-ebiten.KeyNumpad1) + 1)
		}
	}

	// Backspace / Delete
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) || inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
		sr, sc := g.board.Selected[0], g.board.Selected[1]
		if sr >= 0 && sc >= 0 {
			g.board.ClearDigit(sr, sc)
		}
	}

	// R to restart same difficulty
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.restart()
	}

	// Escape to go back to menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.screen = ScreenMenu
	}

	return nil
}

func (g *Game) handleLeftClick(mx, my int) {
	// Reset button
	if render.ResetButtonHit(mx, my) {
		g.restart()
		return
	}

	// Numpad
	digit := render.NumPadDigitAt(mx, my)
	if digit > 0 {
		g.placeDigitOnSelected(digit)
		return
	}

	// Grid
	r, c := render.CellAt(mx, my)
	if r < 0 || c < 0 {
		return
	}

	cell := &g.board.Cells[r][c]
	switch cell.State {
	case board.CellHidden:
		if !g.started {
			g.started = true
			g.startTime = time.Now()
		}
		g.board.Reveal(r, c)
		if g.board.Cells[r][c].State == board.CellRevealed {
			g.board.SelectCell(r, c)
		}
	case board.CellRevealed:
		g.board.SelectCell(r, c)
	case board.CellGiven:
		g.board.SelectCell(r, c)
	}
}

func (g *Game) handleRightClick(mx, my int) {
	r, c := render.CellAt(mx, my)
	if r < 0 || c < 0 {
		return
	}
	g.board.ToggleFlag(r, c)
}

func (g *Game) placeDigitOnSelected(digit int) {
	sr, sc := g.board.Selected[0], g.board.Selected[1]
	if sr < 0 || sc < 0 {
		return
	}
	cell := &g.board.Cells[sr][sc]
	if cell.State != board.CellRevealed {
		return
	}
	if cell.PlayerGuess == digit {
		g.board.ClearDigit(sr, sc)
	} else {
		g.board.PlaceDigit(sr, sc, digit)
	}
}

// --- SCORES ---

func (g *Game) updateScores() error {
	mx, my := ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if render.ScoresBackButtonHit(mx, my) {
			g.screen = ScreenMenu
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		g.screen = ScreenMenu
	}

	return nil
}

// --- DRAW ---

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.screen {
	case ScreenMenu:
		render.DrawMenu(screen, g.menuHover)
	case ScreenPlaying:
		render.DrawGame(screen, g.board, g.hoverR, g.hoverC, g.elapsed, g.difficulty.Name)
	case ScreenScores:
		render.DrawScores(screen)
	}
}

func main() {
	g := NewGame()
	w, h := render.ScreenSize()

	ebiten.SetWindowTitle("Ninesweeper")
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetTPS(30)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
