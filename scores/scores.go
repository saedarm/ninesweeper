package scores

import (
	"encoding/json"
	"os"
	"sort"
)

const maxScoresPerDifficulty = 5
const scoresFile = "ninesweeper_scores.json"

// Entry is a single high score record
type Entry struct {
	Difficulty string `json:"difficulty"`
	Time       int    `json:"time"`
}

// Table holds all high scores
type Table struct {
	Entries []Entry `json:"entries"`
}

var current *Table

func init() {
	current = Load()
}

// Load reads scores from disk (or returns empty table)
func Load() *Table {
	t := &Table{Entries: []Entry{}}
	data, err := os.ReadFile(scoresFile)
	if err != nil {
		return t
	}
	_ = json.Unmarshal(data, t)
	return t
}

// Save writes scores to disk
func Save() {
	data, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(scoresFile, data, 0644)
}

// Add records a new score and returns true if it's a new high score
func Add(difficulty string, timeSec int) bool {
	entry := Entry{Difficulty: difficulty, Time: timeSec}
	current.Entries = append(current.Entries, entry)
	Save()
	return IsHighScore(difficulty, timeSec)
}

// IsHighScore checks if a time would rank in the top scores for a difficulty
func IsHighScore(difficulty string, timeSec int) bool {
	filtered := GetForDifficulty(difficulty)
	if len(filtered) < maxScoresPerDifficulty {
		return true
	}
	return timeSec < filtered[len(filtered)-1].Time
}

// GetForDifficulty returns the top scores for a difficulty, sorted best-first
func GetForDifficulty(difficulty string) []Entry {
	var filtered []Entry
	for _, e := range current.Entries {
		if e.Difficulty == difficulty {
			filtered = append(filtered, e)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Time < filtered[j].Time
	})
	if len(filtered) > maxScoresPerDifficulty {
		filtered = filtered[:maxScoresPerDifficulty]
	}
	return filtered
}

// GetAll returns the full table
func GetAll() *Table {
	return current
}

// BestTime returns the best time for a difficulty, or -1 if none
func BestTime(difficulty string) int {
	entries := GetForDifficulty(difficulty)
	if len(entries) == 0 {
		return -1
	}
	return entries[0].Time
}
