package mapper

import (
	"strconv"

	gqlmodel "github.com/haranton/go-graphql-blog/graph/model"
	"github.com/haranton/go-graphql-blog/internals/models"
)

func MapCommentDomainToGraph(c *models.Comment, all []models.Comment) *gqlmodel.Comment {
	if c == nil {
		return nil
	}

	// конвертируем базовые поля
	var parentID *string
	if c.ParentID != nil {
		s := strconv.Itoa(*c.ParentID)
		parentID = &s
	}

	// собираем ответы (replies) — те комментарии, у которых ParentID == c.ID
	var replies []*gqlmodel.Comment
	for i := range all {
		if all[i].ParentID != nil && *all[i].ParentID == c.ID {
			replies = append(replies, MapCommentDomainToGraph(&all[i], all))
		}
	}

	return &gqlmodel.Comment{
		ID:       strconv.Itoa(c.ID),
		PostID:   strconv.Itoa(c.PostID),
		ParentID: parentID,
		Content:  c.Content,
		Replies:  replies,
	}
}

func MapPostWithCommentsDomainToGraph(p *models.PostWithComments) *gqlmodel.Post {
	if p == nil {
		return nil
	}

	var comments []*gqlmodel.Comment
	for i := range p.Comments {
		if p.Comments[i].ParentID == nil {
			comments = append(comments, MapCommentDomainToGraph(&p.Comments[i], p.Comments))
		}
	}

	return &gqlmodel.Post{
		ID:            strconv.Itoa(p.ID),
		Title:         p.Title,
		Content:       p.Content,
		AllowComments: p.AllowComments,
		Comments:      comments,
	}
}
