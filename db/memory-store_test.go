package db

import (
	"PostCommentService/graph/model"
	"testing"
)

func TestGetPostsMemory(t *testing.T) {
	store := NewMemoryStore()

	post1 := &model.Post{ID: 1, Title: "Post 1"}
	post2 := &model.Post{ID: 2, Title: "Post 2"}

	store.posts[1] = post1
	store.posts[2] = post2

	posts, err := store.GetPosts()
	if err != nil {
		t.Errorf("error was not expected while getting posts: %s", err)
	}

	if len(posts) != 2 {
		t.Errorf("expected 2 posts, got %d", len(posts))
	}
}

func TestGetPostMemory(t *testing.T) {
	store := NewMemoryStore()

	post := &model.Post{ID: 1, Title: "Post 1"}
	store.posts[1] = post

	gotPost, err := store.GetPost(1, 0, 10)
	if err != nil {
		t.Errorf("error was not expected while getting post: %s", err)
	}

	if gotPost == nil {
		t.Errorf("post should not be nil")
		return
	}

	if gotPost.ID != 1 || gotPost.Title != "Post 1" {
		t.Errorf("unexpected values in post: %+v", gotPost)
	}
}

func TestGetCommentsMemory(t *testing.T) {
	store := NewMemoryStore()

	comment1 := &model.Comment{ID: 1, PostID: 1, Author: "Author 1", Content: "Content 1"}
	comment2 := &model.Comment{ID: 2, PostID: 1, Author: "Author 2", Content: "Content 2"}

	store.comments[1] = comment1
	store.comments[2] = comment2

	comments, err := store.GetComments(1, 0, 10)
	if err != nil {
		t.Errorf("error was not expected while getting comments: %s", err)
	}

	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}
}

func TestGetCommentMemory(t *testing.T) {
	store := NewMemoryStore()

	comment := &model.Comment{ID: 1, PostID: 1, Author: "Author", Content: "Content"}
	store.comments[1] = comment

	gotComment, err := store.GetComment(1)
	if err != nil {
		t.Errorf("error was not expected while getting comment: %s", err)
	}

	if gotComment == nil {
		t.Errorf("comment should not be nil")
		return
	}

	if gotComment.ID != 1 || gotComment.Author != "Author" || gotComment.Content != "Content" {
		t.Errorf("unexpected values in comment: %+v", gotComment)
	}
}

func TestCreatePostMemory(t *testing.T) {
	store := NewMemoryStore()

	post, err := store.CreatePost("Title", "Content", "Author")
	if err != nil {
		t.Errorf("error was not expected while creating post: %s", err)
	}

	if post == nil {
		t.Errorf("post should not be nil")
		return
	}

	if post.ID != 1 || post.Title != "Title" || post.Content != "Content" || post.Author != "Author" {
		t.Errorf("unexpected values in post: %+v", post)
	}
}

func TestCreateCommentMemory(t *testing.T) {
	store := NewMemoryStore()

	post, _ := store.CreatePost("Title", "Content", "Author")
	comment, err := store.CreateComment(post.ID, "Author", "Content", nil)
	if err != nil {
		t.Errorf("error was not expected while creating comment: %s", err)
	}

	if comment == nil {
		t.Errorf("comment should not be nil")
		return
	}

	if comment.ID != 1 || comment.PostID != post.ID || comment.Author != "Author" || comment.Content != "Content" {
		t.Errorf("unexpected values in comment: %+v", comment)
	}
}

func TestUpdatePostMemory(t *testing.T) {
	store := NewMemoryStore()

	post, _ := store.CreatePost("Title", "Content", "Author")
	updatedPost, err := store.UpdatePost(post.ID, "Updated Title", "Updated Content")
	if err != nil {
		t.Errorf("error was not expected while updating post: %s", err)
	}

	if updatedPost == nil {
		t.Errorf("updated post should not be nil")
		return
	}

	if updatedPost.ID != post.ID || updatedPost.Title != "Updated Title" || updatedPost.Content != "Updated Content" {
		t.Errorf("unexpected values in updated post: %+v", updatedPost)
	}
}

func TestUpdateCommentMemory(t *testing.T) {
	store := NewMemoryStore()

	post, _ := store.CreatePost("Title", "Content", "Author")
	comment, _ := store.CreateComment(post.ID, "Author", "Content", nil)
	updatedComment, err := store.UpdateComment(comment.ID, "Updated Content")
	if err != nil {
		t.Errorf("error was not expected while updating comment: %s", err)
	}

	if updatedComment == nil {
		t.Errorf("updated comment should not be nil")
		return
	}

	if updatedComment.ID != comment.ID || updatedComment.Content != "Updated Content" {
		t.Errorf("unexpected values in updated comment: %+v", updatedComment)
	}
}

func TestDisableCommentsMemory(t *testing.T) {
	store := NewMemoryStore()

	post, _ := store.CreatePost("Title", "Content", "Author")
	err := store.DisableComments(post.ID)
	if err != nil {
		t.Errorf("error was not expected while disabling comments: %s", err)
	}

	updatedPost, _ := store.GetPost(post.ID, 0, 10)
	if updatedPost.CommentsEnabled {
		t.Errorf("comments should be disabled for the post")
	}
}
