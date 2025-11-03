package memory

import (
	"context"
	"sort"
	"time"

	"github.com/haranton/go-graphql-blog/internals/models"
)

func (st *MemoryStorage) Comment(ctx context.Context, id int) (*models.Comment, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if c, ok := st.commentByID[id]; ok {
		copy := *c
		return &copy, nil
	}

	for _, c := range st.comments {
		if c.ID == id {
			copy := *c
			return &copy, nil
		}
	}

	return nil, nil
}

func (st *MemoryStorage) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.nextCommentID++
	comment.ID = st.nextCommentID
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	copy := *comment
	st.comments = append(st.comments, &copy)
	st.commentByID[copy.ID] = &copy

	st.commentsByPostID[copy.PostID] = append(st.commentsByPostID[copy.PostID], &copy)

	if copy.ParentID != nil {
		st.repliesByParentID[*copy.ParentID] = append(st.repliesByParentID[*copy.ParentID], &copy)
	}

	return &copy, nil
}

func (st *MemoryStorage) ListComments(ctx context.Context, postID, limit, offset int) ([]models.Comment, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	all := st.commentsByPostID[postID]
	if len(all) == 0 {
		return []models.Comment{}, nil
	}

	var topLevel []models.Comment
	for _, c := range all {
		if c.ParentID == nil {
			topLevel = append(topLevel, *c)
		}
	}

	sort.Slice(topLevel, func(i, j int) bool {
		return topLevel[i].CreatedAt.Before(topLevel[j].CreatedAt)
	})

	if offset >= len(topLevel) {
		return []models.Comment{}, nil
	}

	end := offset + limit
	if end > len(topLevel) {
		end = len(topLevel)
	}

	return topLevel[offset:end], nil
}

// Список ответов по parent_id
func (st *MemoryStorage) ListReplies(ctx context.Context, parentID, limit, offset int) ([]models.Comment, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	all := st.repliesByParentID[parentID]
	if len(all) == 0 {
		return []models.Comment{}, nil
	}

	replies := make([]models.Comment, len(all))
	for i, c := range all {
		replies[i] = *c
	}

	sort.Slice(replies, func(i, j int) bool {
		return replies[i].CreatedAt.Before(replies[j].CreatedAt)
	})

	if offset >= len(replies) {
		return []models.Comment{}, nil
	}

	end := offset + limit
	if end > len(replies) {
		end = len(replies)
	}

	return replies[offset:end], nil
}

// Для DataLoader — получить все ответы по набору parentIDs
func (st *MemoryStorage) ListRepliesBatch(ctx context.Context, parentIDs []int) ([]models.Comment, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var results []models.Comment
	for _, id := range parentIDs {
		if replies, ok := st.repliesByParentID[id]; ok {
			for _, c := range replies {
				results = append(results, *c)
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.Before(results[j].CreatedAt)
	})

	return results, nil
}
