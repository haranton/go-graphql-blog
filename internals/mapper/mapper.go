package mapper

import (
	"strconv"

	gqlmodel "github.com/haranton/go-graphql-blog/graph/model"
	"github.com/haranton/go-graphql-blog/internals/models"
)

func MapCommentDomainToGraph(c *models.Comment) *gqlmodel.Comment {
	if c == nil {
		return nil
	}

	var parentID *string
	if c.ParentID != nil {
		s := strconv.Itoa(*c.ParentID)
		parentID = &s
	}

	return &gqlmodel.Comment{
		ID:       strconv.Itoa(c.ID),
		PostID:   strconv.Itoa(c.PostID),
		ParentID: parentID,
		Content:  c.Content,
		Replies:  nil,
	}
}

func MapPostWithCommentsDomainToGraph(p *models.Post) *gqlmodel.Post {
	if p == nil {
		return nil
	}

	return &gqlmodel.Post{
		ID:            strconv.Itoa(p.ID),
		Title:         p.Title,
		Content:       p.Content,
		AllowComments: p.AllowComments,
	}
}
