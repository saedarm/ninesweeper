package sudoku

import (
	"math/rand"
)

const Size = 9
const BoxSize = 3

// Generate creates a fully solved 9x9 Sudoku grid
func Generate(rng *rand.Rand) [Size][Size]int {
	var grid [Size][Size]int
	fillGrid(rng, &grid, 0, 0)
	return grid
}

func fillGrid(rng *rand.Rand, grid *[Size][Size]int, row, col int) bool {
	if row == Size {
		return true
	}
	nextRow, nextCol := row, col+1
	if nextCol == Size {
		nextRow = row + 1
		nextCol = 0
	}

	nums := rng.Perm(Size)
	for _, n := range nums {
		val := n + 1
		if isValid(grid, row, col, val) {
			grid[row][col] = val
			if fillGrid(rng, grid, nextRow, nextCol) {
				return true
			}
			grid[row][col] = 0
		}
	}
	return false
}

func isValid(grid *[Size][Size]int, row, col, val int) bool {
	// Check row
	for c := 0; c < Size; c++ {
		if grid[row][c] == val {
			return false
		}
	}
	// Check column
	for r := 0; r < Size; r++ {
		if grid[r][col] == val {
			return false
		}
	}
	// Check 3x3 box
	boxR, boxC := (row/BoxSize)*BoxSize, (col/BoxSize)*BoxSize
	for r := boxR; r < boxR+BoxSize; r++ {
		for c := boxC; c < boxC+BoxSize; c++ {
			if grid[r][c] == val {
				return false
			}
		}
	}
	return true
}

// IsValidPlacement checks if placing val at (row, col) is valid
// considering only the currently placed values (ignoring the cell itself)
func IsValidPlacement(grid *[Size][Size]int, row, col, val int) bool {
	for c := 0; c < Size; c++ {
		if c != col && grid[row][c] == val {
			return false
		}
	}
	for r := 0; r < Size; r++ {
		if r != row && grid[r][col] == val {
			return false
		}
	}
	boxR, boxC := (row/BoxSize)*BoxSize, (col/BoxSize)*BoxSize
	for r := boxR; r < boxR+BoxSize; r++ {
		for c := boxC; c < boxC+BoxSize; c++ {
			if (r != row || c != col) && grid[r][c] == val {
				return false
			}
		}
	}
	return true
}
