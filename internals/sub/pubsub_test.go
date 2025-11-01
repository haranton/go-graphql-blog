package sub

import (
	"testing"
	"time"

	"github.com/haranton/go-graphql-blog/graph/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentPubSub_Subscribe_Publish(t *testing.T) {
	ps := NewCommentPubSub()
	postID := "1"

	ch := ps.Subscribe(postID)

	comment := &model.Comment{
		ID:       "1",
		PostID:   postID,
		Content:  "Test comment",
		ParentID: nil,
	}

	ps.Publish(comment)

	select {
	case received := <-ch:
		require.NotNil(t, received)
		assert.Equal(t, comment.ID, received.ID)
		assert.Equal(t, comment.Content, received.Content)
	case <-time.After(1 * time.Second):
		t.Fatal("Did not receive comment within timeout")
	}
}

func TestCommentPubSub_MultipleSubscribers(t *testing.T) {
	ps := NewCommentPubSub()
	postID := "1"

	ch1 := ps.Subscribe(postID)
	ch2 := ps.Subscribe(postID)
	ch3 := ps.Subscribe(postID)

	comment := &model.Comment{
		ID:      "1",
		PostID:  postID,
		Content: "Test comment",
	}

	ps.Publish(comment)

	select {
	case received := <-ch1:
		assert.NotNil(t, received)
		assert.Equal(t, comment.ID, received.ID)
	case <-time.After(1 * time.Second):
		t.Fatal("ch1 did not receive comment")
	}

	select {
	case received := <-ch2:
		assert.NotNil(t, received)
		assert.Equal(t, comment.ID, received.ID)
	case <-time.After(1 * time.Second):
		t.Fatal("ch2 did not receive comment")
	}

	select {
	case received := <-ch3:
		assert.NotNil(t, received)
		assert.Equal(t, comment.ID, received.ID)
	case <-time.After(1 * time.Second):
		t.Fatal("ch3 did not receive comment")
	}
}

func TestCommentPubSub_DifferentPosts(t *testing.T) {
	ps := NewCommentPubSub()
	postID1 := "1"
	postID2 := "2"

	ch1 := ps.Subscribe(postID1)
	ch2 := ps.Subscribe(postID2)

	comment1 := &model.Comment{
		ID:      "1",
		PostID:  postID1,
		Content: "Comment for post 1",
	}

	comment2 := &model.Comment{
		ID:      "2",
		PostID:  postID2,
		Content: "Comment for post 2",
	}

	ps.Publish(comment1)
	ps.Publish(comment2)

	select {
	case received := <-ch1:
		assert.Equal(t, comment1.ID, received.ID)
		assert.Equal(t, postID1, received.PostID)
		assert.Equal(t, "Comment for post 1", received.Content)
	case <-time.After(1 * time.Second):
		t.Fatal("ch1 did not receive comment")
	}

	select {
	case received := <-ch2:
		assert.Equal(t, comment2.ID, received.ID)
		assert.Equal(t, postID2, received.PostID)
		assert.Equal(t, "Comment for post 2", received.Content)
	case <-time.After(1 * time.Second):
		t.Fatal("ch2 did not receive comment")
	}

	select {
	case <-ch1:
		t.Fatal("ch1 should not have received comment2")
	default:

	}
}

func TestCommentPubSub_Unsubscribe(t *testing.T) {
	ps := NewCommentPubSub()
	postID := "1"

	ch1 := ps.Subscribe(postID)
	ch2 := ps.Subscribe(postID)

	ps.Unsubscribe(postID, ch1)

	_, ok := <-ch1
	assert.False(t, ok, "Channel should be closed")

	comment := &model.Comment{
		ID:      "1",
		PostID:  postID,
		Content: "Test comment",
	}

	ps.Publish(comment)

	select {
	case received := <-ch2:
		assert.NotNil(t, received)
		assert.Equal(t, comment.ID, received.ID)
	case <-time.After(1 * time.Second):
		t.Fatal("ch2 did not receive comment")
	}

	_, ok = <-ch1
	assert.False(t, ok, "Channel should still be closed")
}

func TestCommentPubSub_EmptySubscribers(t *testing.T) {
	ps := NewCommentPubSub()

	comment := &model.Comment{
		ID:      "1",
		PostID:  "999",
		Content: "Test comment",
	}

	ps.Publish(comment)

	assert.True(t, true)
}

func TestCommentPubSub_PublishUnsubscribed(t *testing.T) {
	ps := NewCommentPubSub()
	postID := "1"

	ch := ps.Subscribe(postID)

	ps.Unsubscribe(postID, ch)

	comment := &model.Comment{
		ID:      "1",
		PostID:  postID,
		Content: "Test comment",
	}

	ps.Publish(comment)

	_, ok := <-ch
	assert.False(t, ok)
}
