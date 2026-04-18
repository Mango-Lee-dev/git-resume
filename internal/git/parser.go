package git

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/wootaiklee/git-resume/pkg/models"
)

// Parser handles git repository operations
type Parser struct {
	repo   *git.Repository
	filter *Filter
	scorer *Scorer
}

// NewParser creates a new git parser for the given repository path
func NewParser(repoPath string) (*Parser, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	return &Parser{
		repo:   repo,
		filter: NewFilter(),
		scorer: NewScorer(),
	}, nil
}

// GetCommits retrieves commits within the specified date range
func (p *Parser) GetCommits(from, to time.Time) ([]models.Commit, error) {
	ref, err := p.repo.Head()
	if err != nil {
		return nil, err
	}

	commitIter, err := p.repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}

	var commits []models.Commit

	err = commitIter.ForEach(func(c *object.Commit) error {
		// Skip commits outside date range
		if c.Author.When.Before(from) {
			return nil
		}
		if c.Author.When.After(to) {
			return nil
		}

		commit := p.parseCommit(c)

		// Apply filtering
		if p.filter.ShouldSkip(commit) {
			return nil
		}

		// Calculate importance score
		commit.Score = p.scorer.Calculate(commit)

		commits = append(commits, commit)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return commits, nil
}

// GetCommitsByMonth retrieves commits for a specific month
func (p *Parser) GetCommitsByMonth(year int, month time.Month) ([]models.Commit, error) {
	from := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	to := from.AddDate(0, 1, 0).Add(-time.Second)
	return p.GetCommits(from, to)
}

// parseCommit converts a go-git commit to our model
func (p *Parser) parseCommit(c *object.Commit) models.Commit {
	files, additions, deletions := p.getCommitStats(c)

	return models.Commit{
		Hash:      c.Hash.String(),
		Author:    c.Author.Name,
		Email:     c.Author.Email,
		Date:      c.Author.When,
		Message:   c.Message,
		Files:     files,
		Additions: additions,
		Deletions: deletions,
	}
}

// getCommitStats calculates file changes for a commit
func (p *Parser) getCommitStats(c *object.Commit) ([]string, int, int) {
	var files []string
	var additions, deletions int

	stats, err := c.Stats()
	if err != nil {
		return files, 0, 0
	}

	for _, stat := range stats {
		// Skip filtered files
		if p.filter.ShouldSkipFile(stat.Name) {
			continue
		}
		files = append(files, stat.Name)
		additions += stat.Addition
		deletions += stat.Deletion
	}

	return files, additions, deletions
}
