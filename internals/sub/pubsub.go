package sub

import (
	"fmt"
	"sync"

	"github.com/haranton/go-graphql-blog/graph/model"
)

type CommentPubSub struct {
	subscribers map[string]map[chan *model.Comment]bool
	mu          sync.RWMutex
}

func NewCommentPubSub() *CommentPubSub {
	return &CommentPubSub{
		subscribers: make(map[string]map[chan *model.Comment]bool),
		mu:          sync.RWMutex{},
	}
}

func (comSub *CommentPubSub) Subscribe(postId string) chan *model.Comment {
	chanClient := make(chan *model.Comment, 1)
	comSub.mu.Lock()
	if comSub.subscribers[postId] == nil {
		comSub.subscribers[postId] = make(map[chan *model.Comment]bool)
	}
	comSub.subscribers[postId][chanClient] = true
	comSub.mu.Unlock()

	return chanClient
}

func (comSub *CommentPubSub) Unsubscribe(postId string, ch chan *model.Comment) {
	comSub.mu.Lock()
	defer comSub.mu.Unlock()

	if subs, ok := comSub.subscribers[postId]; ok {
		if _, exists := subs[ch]; exists {
			delete(subs, ch)
			close(ch)
		}
		if len(subs) == 0 {
			delete(comSub.subscribers, postId)
		}
	}
}

func (comSub *CommentPubSub) Publish(comment *model.Comment) {
	comSub.mu.RLock()
	defer comSub.mu.RUnlock()

	fmt.Printf("Publishing comment to subscribers of postID %s\n", comment.PostID)

	if chSub, ok := comSub.subscribers[comment.PostID]; ok {
		for ch := range chSub {
			select {
			case ch <- comment:
			default:

			}
		}
	}
}
