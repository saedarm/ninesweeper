package board

import (
	"math/rand"

	"github.com/saedarm/ninesweeper/sudoku"
)

const Size = 9

type CellState int

const (
	CellHidden   CellState = iota // Unrevealed
	CellRevealed                  // Revealed (safe, shows sudoku digit + mine count)
	CellFlagged                   // Player flagged as mine
	CellGiven                     // Pre-revealed given (safe, locked)
	CellExploded                  // Player hit a mine
)

type GameState int

const (
	StateReady   GameState = iota // Waiting for first click
	StatePlaying                  // In progress
	StateWon                      // All non-mine cells correctly filled
	StateLost                     // Clicked a mine
)

type Cell struct {
	State       CellState
	IsMine      bool
	SudokuValue int // The correct sudoku digit (1-9)
	PlayerGuess int // What the player has entered (0 = empty)
	MineCount   int // Adjacent mine count (minesweeper number)
	Selected    bool
}

type Board struct {
	Cells     [Size][Size]Cell
	Solution  [Size][Size]int
	GameState GameState
	MineCount int
	FlagCount int
	Selected  [2]int // row, col of selected cell (-1 = none)
}

// New creates a new board with a valid sudoku, placed mines, and pre-revealed givens
func New(rng *rand.Rand, numMines int, numGivens int) *Board {
	b := &Board{
		MineCount: numMines,
		Selected:  [2]int{-1, -1},
	}

	// Step 1: Generate a complete sudoku solution
	b.Solution = sudoku.Generate(rng)

	// Step 2: Place mines randomly
	positions := rng.Perm(Size * Size)
	minesPlaced := 0
	for _, pos := range positions {
		if minesPlaced >= numMines {
			break
		}
		r, c := pos/Size, pos%Size
		b.Cells[r][c].IsMine = true
		minesPlaced++
	}

	// Step 3: Compute mine adjacency counts for all cells
	for r := 0; r < Size; r++ {
		for c := 0; c < Size; c++ {
			b.Cells[r][c].SudokuValue = b.Solution[r][c]
			if !b.Cells[r][c].IsMine {
				count := 0
				b.forEachNeighbor(r, c, func(nr, nc int) {
					if b.Cells[nr][nc].IsMine {
						count++
					}
				})
				b.Cells[r][c].MineCount = count
			}
		}
	}

	// Step 4: Select givens from non-mine cells
	// Spread them across boxes for better gameplay
	nonMineCells := []int{}
	for _, pos := range positions {
		r, c := pos/Size, pos%Size
		if !b.Cells[r][c].IsMine {
			nonMineCells = append(nonMineCells, pos)
		}
	}

	givensPlaced := 0
	for _, pos := range nonMineCells {
		if givensPlaced >= numGivens {
			break
		}
		r, c := pos/Size, pos%Size
		b.Cells[r][c].State = CellGiven
		b.Cells[r][c].PlayerGuess = b.Cells[r][c].SudokuValue
		givensPlaced++
	}

	b.GameState = StateReady
	return b
}

func (b *Board) forEachNeighbor(row, col int, fn func(r, c int)) {
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			nr, nc := row+dr, col+dc
			if nr >= 0 && nr < Size && nc >= 0 && nc < Size {
				fn(nr, nc)
			}
		}
	}
}

// Reveal attempts to reveal a hidden cell
func (b *Board) Reveal(row, col int) {
	if b.GameState == StateWon || b.GameState == StateLost {
		return
	}
	if b.GameState == StateReady {
		// First click: ensure it's not a mine
		b.ensureSafeFirstClick(row, col)
		b.GameState = StatePlaying
	}

	cell := &b.Cells[row][col]
	if cell.State != CellHidden {
		return
	}

	if cell.IsMine {
		cell.State = CellExploded
		b.GameState = StateLost
		return
	}

	cell.State = CellRevealed

	// Auto-reveal (flood fill) if mine count is 0 AND sudoku value is shown
	if cell.MineCount == 0 {
		b.forEachNeighbor(row, col, func(nr, nc int) {
			if b.Cells[nr][nc].State == CellHidden && !b.Cells[nr][nc].IsMine {
				b.Reveal(nr, nc)
			}
		})
	}
}

// ensureSafeFirstClick moves a mine away from the first clicked cell
func (b *Board) ensureSafeFirstClick(row, col int) {
	cell := &b.Cells[row][col]
	if !cell.IsMine {
		return
	}

	// Move the mine to a non-mine, non-given cell
	cell.IsMine = false
	moved := false
	for r := 0; r < Size && !moved; r++ {
		for c := 0; c < Size && !moved; c++ {
			if !b.Cells[r][c].IsMine && (r != row || c != col) && b.Cells[r][c].State != CellGiven {
				b.Cells[r][c].IsMine = true
				moved = true
			}
		}
	}
	// Recompute all adjacency counts
	for r := 0; r < Size; r++ {
		for c := 0; c < Size; c++ {
			if !b.Cells[r][c].IsMine {
				count := 0
				b.forEachNeighbor(r, c, func(nr, nc int) {
					if b.Cells[nr][nc].IsMine {
						count++
					}
				})
				b.Cells[r][c].MineCount = count
			}
		}
	}
}

// ToggleFlag toggles the mine flag on a hidden cell
func (b *Board) ToggleFlag(row, col int) {
	if b.GameState == StateWon || b.GameState == StateLost {
		return
	}
	cell := &b.Cells[row][col]
	if cell.State == CellHidden {
		cell.State = CellFlagged
		b.FlagCount++
	} else if cell.State == CellFlagged {
		cell.State = CellHidden
		b.FlagCount--
	}
}

// PlaceDigit places a sudoku guess on a revealed cell
func (b *Board) PlaceDigit(row, col, digit int) {
	if b.GameState == StateWon || b.GameState == StateLost {
		return
	}
	cell := &b.Cells[row][col]
	if cell.State != CellRevealed {
		return
	}
	if digit < 0 || digit > 9 {
		return
	}
	cell.PlayerGuess = digit
	b.checkWin()
}

// ClearDigit removes the player's guess from a revealed cell
func (b *Board) ClearDigit(row, col int) {
	if b.GameState == StateWon || b.GameState == StateLost {
		return
	}
	cell := &b.Cells[row][col]
	if cell.State != CellRevealed || cell.State == CellGiven {
		return
	}
	cell.PlayerGuess = 0
}

// SelectCell selects a cell (for digit input)
func (b *Board) SelectCell(row, col int) {
	if b.Selected[0] >= 0 && b.Selected[1] >= 0 {
		b.Cells[b.Selected[0]][b.Selected[1]].Selected = false
	}
	b.Selected = [2]int{row, col}
	b.Cells[row][col].Selected = true
}

// ClearSelection clears the current selection
func (b *Board) ClearSelection() {
	if b.Selected[0] >= 0 && b.Selected[1] >= 0 {
		b.Cells[b.Selected[0]][b.Selected[1]].Selected = false
	}
	b.Selected = [2]int{-1, -1}
}

// checkWin checks if all non-mine cells are revealed and correctly filled
func (b *Board) checkWin() {
	for r := 0; r < Size; r++ {
		for c := 0; c < Size; c++ {
			cell := &b.Cells[r][c]
			if cell.IsMine {
				continue
			}
			if cell.State != CellRevealed && cell.State != CellGiven {
				return // Still hidden cells
			}
			if cell.PlayerGuess != cell.SudokuValue {
				return // Wrong or empty guess
			}
		}
	}
	b.GameState = StateWon
}

// IsConflict checks if a digit at (row, col) conflicts with sudoku rules
// based on currently placed guesses/givens
func (b *Board) IsConflict(row, col int) bool {
	cell := &b.Cells[row][col]
	if cell.PlayerGuess == 0 {
		return false
	}
	val := cell.PlayerGuess

	// Check row
	for c := 0; c < Size; c++ {
		if c != col && b.Cells[row][c].PlayerGuess == val {
			return true
		}
	}
	// Check column
	for r := 0; r < Size; r++ {
		if r != row && b.Cells[r][col].PlayerGuess == val {
			return true
		}
	}
	// Check 3x3 box
	boxR, boxC := (row/3)*3, (col/3)*3
	for r := boxR; r < boxR+3; r++ {
		for c := boxC; c < boxC+3; c++ {
			if (r != row || c != col) && b.Cells[r][c].PlayerGuess == val {
				return true
			}
		}
	}
	return false
}

// RemainingMines returns mines minus flags
func (b *Board) RemainingMines() int {
	return b.MineCount - b.FlagCount
}
