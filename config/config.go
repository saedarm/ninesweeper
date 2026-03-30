package config

// Difficulty holds the parameters for a game difficulty level
type Difficulty struct {
	Name   string
	Mines  int
	Givens int
}

var (
	Easy   = Difficulty{Name: "Easy", Mines: 8, Givens: 25}
	Medium = Difficulty{Name: "Medium", Mines: 12, Givens: 20}
	Hard   = Difficulty{Name: "Hard", Mines: 15, Givens: 15}
	Expert = Difficulty{Name: "Expert", Mines: 18, Givens: 12}

	All = []Difficulty{Easy, Medium, Hard, Expert}
)
