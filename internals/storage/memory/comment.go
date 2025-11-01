package memory

import (
	"context"
	"time"

	"github.com/haranton/go-graphql-blog/internals/models"
)

func (st *MemoryStorage) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.nextCommentID++
	comment.ID = st.nextCommentID
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	copy := *comment
	st.comments = append(st.comments, &copy)

	return &copy, nil
}

// ListComments — верхний уровень комментариев (parent_id IS NULL)
func (st *MemoryStorage) ListComments(ctx context.Context, postID, limit, offset int) ([]models.Comment, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var filtered []models.Comment
	for _, c := range st.comments {
		if c.PostID == postID && c.ParentID == nil {
			filtered = append(filtered, *c)
		}
	}

	sortByCreatedAt(filtered)

	if offset >= len(filtered) {
		return []models.Comment{}, nil
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], nil
}

// ListReplies — одиночный parent_id (для сервисного слоя)
func (st *MemoryStorage) ListReplies(ctx context.Context, parentID, limit, offset int) ([]models.Comment, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var filtered []models.Comment
	for _, c := range st.comments {
		if c.ParentID != nil && *c.ParentID == parentID {
			filtered = append(filtered, *c)
		}
	}

	sortByCreatedAt(filtered)

	if offset >= len(filtered) {
		return []models.Comment{}, nil
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], nil
}

// ListRepliesBatch — для DataLoader
func (st *MemoryStorage) ListRepliesBatch(ctx context.Context, parentIDs []int) ([]models.Comment, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	ids := make(map[int]struct{}, len(parentIDs))
	for _, id := range parentIDs {
		ids[id] = struct{}{}
	}

	var filtered []models.Comment
	for _, c := range st.comments {
		if c.ParentID != nil {
			if _, ok := ids[*c.ParentID]; ok {
				filtered = append(filtered, *c)
			}
		}
	}

	sortByCreatedAt(filtered)

	return filtered, nil
}
