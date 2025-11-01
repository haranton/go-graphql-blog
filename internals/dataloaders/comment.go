package dataloaders

import (
	"context"
	"strconv"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/haranton/go-graphql-blog/graph/model"
	"github.com/haranton/go-graphql-blog/internals/mapper"
	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/haranton/go-graphql-blog/internals/storage"
)

const timeWaitBatch = 2

type CommentLoader struct {
	Loader *dataloader.Loader
}

func NewCommentLoader(st storage.Storage) *CommentLoader {
	batchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var ids []int
		for _, k := range keys {
			id, _ := strconv.Atoi(k.String())
			ids = append(ids, id)
		}

		comments, err := st.ListRepliesBatch(ctx, ids)
		results := make([]*dataloader.Result, len(keys))

		if err != nil {
			for i := range results {
				results[i] = &dataloader.Result{Error: err}
			}
			return results
		}

		grouped := make(map[int][]models.Comment)
		for _, c := range comments {
			grouped[*c.ParentID] = append(grouped[*c.ParentID], c)
		}

		for i, k := range keys {
			id, _ := strconv.Atoi(k.String())
			results[i] = &dataloader.Result{Data: grouped[id]}
		}

		return results

	}

	loader := dataloader.NewBatchedLoader(batchFn, dataloader.WithWait(timeWaitBatch*time.Millisecond))
	return &CommentLoader{Loader: loader}

}

func (l *CommentLoader) LoadReplies(ctx context.Context, parentID string) ([]*model.Comment, error) {
	thunk := l.Loader.Load(ctx, dataloader.StringKey(parentID))
	res, err := thunk()
	if err != nil {
		return nil, err
	}

	domainComments := res.([]models.Comment)
	mapped := make([]*model.Comment, 0, len(domainComments))
	for _, c := range domainComments {
		mapped = append(mapped, mapper.MapCommentDomainToGraph(&c))
	}
	return mapped, nil
}
