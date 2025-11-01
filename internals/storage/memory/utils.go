package memory

import (
	"sort"

	"github.com/haranton/go-graphql-blog/internals/models"
)

func sortByCreatedAt(comments []models.Comment) {
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})
}
