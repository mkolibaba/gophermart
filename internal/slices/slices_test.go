package slices

import (
	"testing"
	"time"
)

func TestSortByTimeDesc(t *testing.T) {
	type document struct {
		ID       int
		IssuedAt time.Time
	}

	now := time.Now()

	docs := []document{
		{
			ID:       2,
			IssuedAt: now,
		},
		{
			ID:       3,
			IssuedAt: now.Add(-24 * time.Hour),
		},
		{
			ID:       1,
			IssuedAt: now.Add(1 * time.Minute),
		},
	}

	SortByTimeDesc(docs, func(d document) time.Time {
		return d.IssuedAt
	})

	for i := 0; i < len(docs); i++ {
		if docs[i].ID != i+1 {
			t.Errorf("want docs with id = %d at place %d, got %d", docs[i].ID, docs[i].ID-1, i)
		}
	}
}
