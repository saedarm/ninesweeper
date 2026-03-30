package render

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/saedarm/ninesweeper/board"
	"github.com/saedarm/ninesweeper/config"
	"github.com/saedarm/ninesweeper/scores"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

const (
	CellSize    = 56
	GridPadding = 16
	HeaderH     = 44
	FooterH     = 52
	BoxLineW    = 3.0
	CellLineW   = 1.0
)

// Font faces at different sizes — initialized once
var (
	fontSmall     font.Face // 12px - mine counts, labels
	fontMedium    font.Face // 16px - header, buttons, instructions
	fontLarge     font.Face // 28px - big sudoku digits
	fontTitleBold font.Face // 36px bold - menu title
)

func init() {
	ttRegular, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	ttBold, err := opentype.Parse(gobold.TTF)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72

	fontSmall, err = opentype.NewFace(ttRegular, &opentype.FaceOptions{Size: 12, DPI: dpi, Hinting: font.HintingFull})
	if err != nil {
		log.Fatal(err)
	}
	fontMedium, err = opentype.NewFace(ttRegular, &opentype.FaceOptions{Size: 16, DPI: dpi, Hinting: font.HintingFull})
	if err != nil {
		log.Fatal(err)
	}
	fontLarge, err = opentype.NewFace(ttBold, &opentype.FaceOptions{Size: 28, DPI: dpi, Hinting: font.HintingFull})
	if err != nil {
		log.Fatal(err)
	}
	fontTitleBold, err = opentype.NewFace(ttBold, &opentype.FaceOptions{Size: 36, DPI: dpi, Hinting: font.HintingFull})
	if err != nil {
		log.Fatal(err)
	}
}

var (
	colorBg          = color.RGBA{R: 0xF5, G: 0xF5, B: 0xF0, A: 0xFF}
	colorGridLine    = color.RGBA{R: 0xBB, G: 0xBB, B: 0xBB, A: 0xFF}
	colorBoxLine     = color.RGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xFF}
	colorHidden      = color.RGBA{R: 0xB0, G: 0xBE, B: 0xC5, A: 0xFF}
	colorHiddenHover = color.RGBA{R: 0x9E, G: 0xAE, B: 0xB8, A: 0xFF}
	colorRevealed    = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFE, A: 0xFF}
	colorGiven       = color.RGBA{R: 0xE8, G: 0xF5, B: 0xE9, A: 0xFF}
	colorSelected    = color.RGBA{R: 0xFF, G: 0xF9, B: 0xC4, A: 0xFF}
	colorFlagged     = color.RGBA{R: 0xFF, G: 0xCC, B: 0xBC, A: 0xFF}
	colorExploded    = color.RGBA{R: 0xEF, G: 0x53, B: 0x50, A: 0xFF}
	colorHeaderBg    = color.RGBA{R: 0xE0, G: 0xE0, B: 0xDD, A: 0xFF}
	colorBtnBg       = color.RGBA{R: 0xFF, G: 0xEB, B: 0x3B, A: 0xFF}
	colorBtnHover    = color.RGBA{R: 0xFF, G: 0xF1, B: 0x76, A: 0xFF}
	colorNumPadBg    = color.RGBA{R: 0xE8, G: 0xE8, B: 0xE8, A: 0xFF}
	colorWinBanner   = color.RGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xE0}
	colorLostBanner  = color.RGBA{R: 0xEF, G: 0x44, B: 0x44, A: 0xE0}
	colorTitle       = color.RGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xFF}
	colorSubtle      = color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xFF}
	colorAccent      = color.RGBA{R: 0x15, G: 0x65, B: 0xC0, A: 0xFF}
	colorWhite       = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	colorGivenDigit  = color.RGBA{R: 0x2E, G: 0x7D, B: 0x32, A: 0xFF}
	colorPlayerDigit = color.RGBA{R: 0x15, G: 0x65, B: 0xC0, A: 0xFF}
	colorConflict    = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	colorFlagText    = color.RGBA{R: 0xBF, G: 0x36, B: 0x0C, A: 0xFF}
	colorMineText    = color.RGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xFF}

	mineCountColors = [9]color.Color{
		color.Transparent,
		color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}, // 1 blue
		color.RGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xFF}, // 2 green
		color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, // 3 red
		color.RGBA{R: 0x00, G: 0x00, B: 0x80, A: 0xFF}, // 4 navy
		color.RGBA{R: 0x80, G: 0x00, B: 0x00, A: 0xFF}, // 5 maroon
		color.RGBA{R: 0x00, G: 0x80, B: 0x80, A: 0xFF}, // 6 teal
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}, // 7 black
		color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}, // 8 gray
	}
)

// drawRect draws a filled rectangle. Wraps vector call.
func drawRect(dst *ebiten.Image, x, y, w, h float32, clr color.Color) {
	vector.FillRect(dst, x, y, w, h, clr, true)
}

// strokeRect draws a stroked rectangle.
func strokeRect(dst *ebiten.Image, x, y, w, h, strokeW float32, clr color.Color) {
	vector.StrokeRect(dst, x, y, w, h, strokeW, clr, true)
}

// strokeLine draws a line.
func strokeLine(dst *ebiten.Image, x0, y0, x1, y1, strokeW float32, clr color.Color) {
	vector.StrokeLine(dst, x0, y0, x1, y1, strokeW, clr, true)
}

// ScreenSize returns the window width and height.
func ScreenSize() (int, int) {
	w := GridPadding*2 + CellSize*board.Size
	h := HeaderH + GridPadding + CellSize*board.Size + GridPadding + FooterH
	return w, h
}

// GridOrigin returns top-left pixel of the 9x9 grid.
func GridOrigin() (int, int) {
	return GridPadding, HeaderH + GridPadding
}

// CellAt converts screen coordinates to grid (row, col). Returns -1,-1 if outside.
func CellAt(x, y int) (int, int) {
	ox, oy := GridOrigin()
	if x < ox || y < oy {
		return -1, -1
	}
	col := (x - ox) / CellSize
	row := (y - oy) / CellSize
	if row < 0 || row >= board.Size || col < 0 || col >= board.Size {
		return -1, -1
	}
	return row, col
}

// NumPadDigitAt returns 1-9 if the point hits a numpad button, else 0.
func NumPadDigitAt(x, y int) int {
	sw, sh := ScreenSize()
	footerY := sh - FooterH
	btnW, btnH := 40, 34
	spacing := 4
	totalW := 9*btnW + 8*spacing
	startX := (sw - totalW) / 2
	btnY := footerY + (FooterH-btnH)/2

	if y < btnY || y >= btnY+btnH {
		return 0
	}
	for d := 1; d <= 9; d++ {
		bx := startX + (d-1)*(btnW+spacing)
		if x >= bx && x < bx+btnW {
			return d
		}
	}
	return 0
}

// ResetButtonHit checks if the point is on the smiley/reset button.
func ResetButtonHit(x, y int) bool {
	sw, _ := ScreenSize()
	bw, bh := 36, 28
	bx := sw/2 - bw/2
	by := (HeaderH - bh) / 2
	return x >= bx && x < bx+bw && y >= by && y < by+bh
}

// --- MENU SCREEN ---

// MenuButton holds layout info for a difficulty button on the menu.
type MenuButton struct {
	X, Y, W, H int
	Difficulty config.Difficulty
}

// MenuButtons returns the layout of the 4 difficulty buttons.
func MenuButtons() []MenuButton {
	sw, sh := ScreenSize()
	btnW, btnH := 260, 44
	spacing := 16
	totalH := len(config.All)*btnH + (len(config.All)-1)*spacing
	startY := sh/2 - totalH/2 + 20
	startX := sw/2 - btnW/2

	buttons := make([]MenuButton, len(config.All))
	for i, d := range config.All {
		buttons[i] = MenuButton{
			X: startX, Y: startY + i*(btnH+spacing),
			W: btnW, H: btnH, Difficulty: d,
		}
	}
	return buttons
}

// MenuButtonAt returns the index of the menu button at (x, y), or -1.
func MenuButtonAt(x, y int) int {
	for i, btn := range MenuButtons() {
		if x >= btn.X && x < btn.X+btn.W && y >= btn.Y && y < btn.Y+btn.H {
			return i
		}
	}
	return -1
}

// DrawMenu renders the title screen with difficulty selector.
func DrawMenu(screen *ebiten.Image, hoverIdx int) {
	sw, _ := ScreenSize()
	screen.Fill(colorBg)

	text.Draw(screen, "NINESWEEPER", fontTitleBold, sw/2-130, 70, colorTitle)
	text.Draw(screen, "Sudoku + Minesweeper", fontSmall, sw/2-68, 92, colorSubtle)

	for i, btn := range MenuButtons() {
		bg := colorBtnBg
		if i == hoverIdx {
			bg = colorBtnHover
		}
		drawRect(screen, float32(btn.X), float32(btn.Y), float32(btn.W), float32(btn.H), bg)
		strokeRect(screen, float32(btn.X), float32(btn.Y), float32(btn.W), float32(btn.H), 2, colorBoxLine)

		label := fmt.Sprintf("%s  (%d mines, %d givens)", btn.Difficulty.Name, btn.Difficulty.Mines, btn.Difficulty.Givens)
		text.Draw(screen, label, fontSmall, btn.X+14, btn.Y+28, colorTitle)

		best := scores.BestTime(btn.Difficulty.Name)
		if best >= 0 {
			bestStr := fmt.Sprintf("Best: %ds", best)
			text.Draw(screen, bestStr, fontSmall, btn.X+btn.W-70, btn.Y+28, colorAccent)
		}
	}

	_, sh := ScreenSize()
	text.Draw(screen, "Left click: reveal    Right click: flag", fontSmall, sw/2-130, sh-60, colorSubtle)
	text.Draw(screen, "1-9: place digit    Backspace: clear", fontSmall, sw/2-120, sh-42, colorSubtle)
	text.Draw(screen, "R: restart    Esc: menu", fontSmall, sw/2-75, sh-24, colorSubtle)
}

// --- GAME SCREEN ---

// DrawGame renders the active game.
func DrawGame(screen *ebiten.Image, b *board.Board, hoverR, hoverC int, elapsed int, difficulty string) {
	screen.Fill(colorBg)
	drawHeader(screen, b, elapsed, difficulty)
	drawCells(screen, b, hoverR, hoverC)
	drawGridLines(screen)
	drawNumPad(screen)
	drawBanner(screen, b, elapsed)
}

func drawHeader(screen *ebiten.Image, b *board.Board, elapsed int, difficulty string) {
	sw, _ := ScreenSize()
	drawRect(screen, 0, 0, float32(sw), float32(HeaderH), colorHeaderBg)

	text.Draw(screen, fmt.Sprintf("MINES: %02d", b.RemainingMines()), fontSmall, 10, 18, colorTitle)
	text.Draw(screen, difficulty, fontSmall, 10, 36, colorSubtle)
	text.Draw(screen, fmt.Sprintf("TIME: %03d", elapsed), fontSmall, sw-95, 28, colorTitle)

	bw, bh := 36, 28
	bx := float32(sw/2 - bw/2)
	by := float32((HeaderH - bh) / 2)
	drawRect(screen, bx, by, float32(bw), float32(bh), colorBtnBg)
	strokeRect(screen, bx, by, float32(bw), float32(bh), 2, colorBoxLine)
	face := ":-)"
	switch b.GameState {
	case board.StateLost:
		face = "X-("
	case board.StateWon:
		face = "B-)"
	}
	text.Draw(screen, face, fontSmall, int(bx)+5, int(by)+20, colorTitle)
}

func drawCells(screen *ebiten.Image, b *board.Board, hoverR, hoverC int) {
	ox, oy := GridOrigin()

	for r := 0; r < board.Size; r++ {
		for c := 0; c < board.Size; c++ {
			cell := &b.Cells[r][c]
			x := float32(ox + c*CellSize)
			y := float32(oy + r*CellSize)

			bg := colorHidden
			switch cell.State {
			case board.CellHidden:
				if r == hoverR && c == hoverC {
					bg = colorHiddenHover
				}
			case board.CellRevealed:
				bg = colorRevealed
				if cell.Selected {
					bg = colorSelected
				}
			case board.CellGiven:
				bg = colorGiven
			case board.CellFlagged:
				bg = colorFlagged
			case board.CellExploded:
				bg = colorExploded
			}
			drawRect(screen, x+1, y+1, float32(CellSize)-2, float32(CellSize)-2, bg)

			ix := int(x)
			iy := int(y)

			switch cell.State {
			case board.CellHidden:
				// nothing
			case board.CellGiven:
				drawBigDigit(screen, cell.SudokuValue, ix, iy, colorGivenDigit)
				if cell.MineCount > 0 {
					drawSmallMineCount(screen, cell.MineCount, ix, iy)
				}
			case board.CellRevealed:
				if cell.MineCount > 0 {
					drawSmallMineCount(screen, cell.MineCount, ix, iy)
				}
				if cell.PlayerGuess > 0 {
					clr := colorPlayerDigit
					if b.IsConflict(r, c) {
						clr = colorConflict
					}
					drawBigDigit(screen, cell.PlayerGuess, ix, iy, clr)
				}
			case board.CellFlagged:
				text.Draw(screen, "F", fontMedium, ix+CellSize/2-5, iy+CellSize/2+6, colorFlagText)
			case board.CellExploded:
				text.Draw(screen, "*", fontLarge, ix+CellSize/2-9, iy+CellSize/2+10, colorWhite)
			}

			if b.GameState == board.StateLost && cell.IsMine &&
				cell.State != board.CellExploded && cell.State != board.CellFlagged {
				text.Draw(screen, "*", fontLarge, ix+CellSize/2-9, iy+CellSize/2+10, colorMineText)
			}
		}
	}
}

func drawGridLines(screen *ebiten.Image) {
	ox, oy := GridOrigin()
	fox := float32(ox)
	foy := float32(oy)
	gridW := float32(board.Size * CellSize)
	gridH := float32(board.Size * CellSize)

	for i := 0; i <= board.Size; i++ {
		fi := float32(i * CellSize)
		strokeLine(screen, fox+fi, foy, fox+fi, foy+gridH, CellLineW, colorGridLine)
		strokeLine(screen, fox, foy+fi, fox+gridW, foy+fi, CellLineW, colorGridLine)
	}
	for i := 0; i <= 3; i++ {
		fi := float32(i * 3 * CellSize)
		strokeLine(screen, fox+fi, foy, fox+fi, foy+gridH, BoxLineW, colorBoxLine)
		strokeLine(screen, fox, foy+fi, fox+gridW, foy+fi, BoxLineW, colorBoxLine)
	}
}

func drawNumPad(screen *ebiten.Image) {
	sw, sh := ScreenSize()
	footerY := sh - FooterH
	drawRect(screen, 0, float32(footerY), float32(sw), float32(FooterH), colorHeaderBg)

	btnW, btnH := 40, 34
	spacing := 4
	totalW := 9*btnW + 8*spacing
	startX := (sw - totalW) / 2
	btnY := footerY + (FooterH-btnH)/2

	for d := 1; d <= 9; d++ {
		bx := startX + (d-1)*(btnW+spacing)
		drawRect(screen, float32(bx), float32(btnY), float32(btnW), float32(btnH), colorNumPadBg)
		strokeRect(screen, float32(bx), float32(btnY), float32(btnW), float32(btnH), 1, colorGridLine)
		text.Draw(screen, fmt.Sprintf("%d", d), fontMedium, bx+13, btnY+24, colorTitle)
	}
}

func drawBanner(screen *ebiten.Image, b *board.Board, elapsed int) {
	if b.GameState != board.StateWon && b.GameState != board.StateLost {
		return
	}
	sw, sh := ScreenSize()
	bannerH := float32(48)
	bannerY := float32(sh)/2 - bannerH/2

	if b.GameState == board.StateWon {
		drawRect(screen, 0, bannerY, float32(sw), bannerH, colorWinBanner)
		text.Draw(screen, fmt.Sprintf("YOU WIN!  Time: %ds", elapsed), fontMedium, 20, int(bannerY)+22, colorWhite)
		text.Draw(screen, "Smiley = restart    Esc = menu", fontSmall, 20, int(bannerY)+40, colorWhite)
	} else {
		drawRect(screen, 0, bannerY, float32(sw), bannerH, colorLostBanner)
		text.Draw(screen, "BOOM!", fontMedium, 20, int(bannerY)+22, colorWhite)
		text.Draw(screen, "Smiley = restart    Esc = menu", fontSmall, 20, int(bannerY)+40, colorWhite)
	}
}

// --- SCORES SCREEN ---

// ScoresBackButtonHit checks if the "Back" button is hit on the scores screen.
func ScoresBackButtonHit(x, y int) bool {
	sw, sh := ScreenSize()
	btnW, btnH := 100, 36
	bx := sw/2 - btnW/2
	by := sh - 70
	return x >= bx && x < bx+btnW && y >= by && y < by+btnH
}

// DrawScores renders the high scores screen.
func DrawScores(screen *ebiten.Image) {
	sw, sh := ScreenSize()
	screen.Fill(colorBg)

	text.Draw(screen, "HIGH SCORES", fontTitleBold, sw/2-105, 55, colorTitle)

	y := 100
	for _, d := range config.All {
		text.Draw(screen, d.Name, fontMedium, 30, y, colorAccent)
		y += 24
		entries := scores.GetForDifficulty(d.Name)
		if len(entries) == 0 {
			text.Draw(screen, "  No scores yet", fontSmall, 30, y, colorSubtle)
			y += 18
		} else {
			for i, e := range entries {
				medal := "  "
				if i == 0 {
					medal = "* "
				}
				text.Draw(screen, fmt.Sprintf("%s%d. %ds", medal, i+1, e.Time), fontSmall, 30, y, colorTitle)
				y += 18
			}
		}
		y += 10
	}

	btnW, btnH := 100, 36
	bx := float32(sw/2 - btnW/2)
	by := float32(sh - 70)
	drawRect(screen, bx, by, float32(btnW), float32(btnH), colorBtnBg)
	strokeRect(screen, bx, by, float32(btnW), float32(btnH), 2, colorBoxLine)
	text.Draw(screen, "BACK", fontMedium, int(bx)+30, int(by)+24, colorTitle)
}

// --- DRAWING HELPERS ---

func drawBigDigit(screen *ebiten.Image, digit int, cx, cy int, clr color.Color) {
	s := fmt.Sprintf("%d", digit)
	text.Draw(screen, s, fontLarge, cx+CellSize/2-9, cy+CellSize/2+10, clr)
}

func drawSmallMineCount(screen *ebiten.Image, count int, cx, cy int) {
	if count <= 0 || count > 8 {
		return
	}
	clr := mineCountColors[count]
	s := fmt.Sprintf("%d", count)
	text.Draw(screen, s, fontSmall, cx+CellSize-16, cy+14, clr)
}
