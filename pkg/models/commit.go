package models

import "time"

// Commit represents a parsed git commit with relevant metadata
type Commit struct {
	Hash      string    `json:"hash"`
	Author    string    `json:"author"`
	Email     string    `json:"email"`
	Date      time.Time `json:"date"`
	Message   string    `json:"message"`
	Files     []string  `json:"files"`
	Additions int       `json:"additions"`
	Deletions int       `json:"deletions"`
	Score     int       `json:"score"` // Importance score (0-100)
}

// Category represents the type of change
type Category string

const (
	CategoryFeature  Category = "Feature"
	CategoryFix      Category = "Fix"
	CategoryRefactor Category = "Refactor"
	CategoryDocs     Category = "Docs"
	CategoryTest     Category = "Test"
	CategoryChore    Category = "Chore"
)

// CommitBatch represents a group of related commits for batch processing
type CommitBatch struct {
	Commits  []Commit
	Project  string
	FromDate time.Time
	ToDate   time.Time
}
