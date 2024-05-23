package db

import (
	"PostCommentService/graph/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	ps := NewPostgresStore(db)

	rows := sqlmock.NewRows([]string{"id", "title", "content", "comments_enabled", "author"}).
		AddRow(1, "Test title 1", "Test content 1", true, "Test author 1").
		AddRow(2, "Test title 2", "Test content 2", false, "Test author 2")

	mock.ExpectQuery("SELECT id, title, content, comments_enabled, author FROM posts").WillReturnRows(rows)

	posts, err := ps.GetPosts()
	if err != nil {
		t.Errorf("error was not expected while getting posts: %s", err)
	}

	if len(posts) != 2 {
		t.Errorf("expected length of posts list to be '2', got '%v'", len(posts))
	}

	post1 := posts[0]
	if post1.ID != 1 || post1.Title != "Test title 1" || post1.Content != "Test content 1" || post1.Author != "Test author 1" || post1.CommentsEnabled != true {
		t.Errorf("unexpected values in post1: %+v", post1)
	}

	post2 := posts[1]
	if post2.ID != 2 || post2.Title != "Test title 2" || post2.Content != "Test content 2" || post2.Author != "Test author 2" || post2.CommentsEnabled != false {
		t.Errorf("unexpected values in post2: %+v", post2)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ps := NewPostgresStore(db)

	rows := sqlmock.NewRows([]string{"id", "title", "content", "comments_enabled", "author"}).
		AddRow(1, "Test title", "Test content", true, "Test author")

	mock.ExpectQuery("^SELECT (.+) FROM posts WHERE id = \\$1$").WithArgs(1).WillReturnRows(rows)

	commentRows := sqlmock.NewRows([]string{"id", "post_id", "author", "content", "parent_id"}).
		AddRow(1, 1, "Comment author", "Comment content", nil)

	mock.ExpectQuery("^SELECT (.+) FROM comments WHERE post_id = \\$1 ORDER BY id ASC").WithArgs(1).WillReturnRows(commentRows)

	post, err := ps.GetPost(1, 0, 10)
	if err != nil {
		t.Errorf("error was not expected while getting post: %s", err)
	}

	if post == nil {
		t.Errorf("post should not be nil")
		return
	}

	if post.ID != 1 || post.Title != "Test title" || post.Content != "Test content" || post.Author != "Test author" || post.CommentsEnabled != true {
		t.Errorf("unexpected values in post: %+v", post)
	}

	if len(post.Comments) != 1 {
		t.Errorf("expected length of comments list to be '1', got '%v'", len(post.Comments))
	}

	comment := post.Comments[0]
	if comment.ID != 1 || comment.PostID != 1 || comment.Author != "Comment author" || comment.Content != "Comment content" {
		t.Errorf("unexpected values in comment: %+v", comment)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ps := NewPostgresStore(db)

	rows := sqlmock.NewRows([]string{"id", "post_id", "author", "content", "parent_id"}).
		AddRow(1, 1, "Comment author", "Comment content", nil)

	mock.ExpectQuery("^SELECT (.+) FROM comments WHERE id = \\$1$").WithArgs(1).WillReturnRows(rows)

	comment, err := ps.GetComment(1)
	if err != nil {
		t.Errorf("error was not expected while getting comment: %s", err)
	}

	if comment == nil {
		t.Errorf("comment should not be nil")
		return
	}

	if comment.ID != 1 || comment.PostID != 1 || comment.Author != "Comment author" || comment.Content != "Comment content" {
		t.Errorf("unexpected values in comment: %+v", comment)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetComments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ps := NewPostgresStore(db)

	commentRows := sqlmock.NewRows([]string{"id", "post_id", "author", "content", "parent_id"}).
		AddRow(1, 1, "Comment author", "Comment content", nil).
		AddRow(2, 1, "Another author", "Another content", nil)

	mock.ExpectQuery("^SELECT (.+) FROM comments WHERE post_id = \\$1 ORDER BY id ASC").WithArgs(1).WillReturnRows(commentRows)

	comments, err := ps.GetComments(1, 0, 10)
	if err != nil {
		t.Errorf("error was not expected while getting comments: %s", err)
	}

	if len(comments) != 2 {
		t.Errorf("expected length of comments list to be '2', got '%v'", len(comments))
	}

	firstComment := comments[0]
	if firstComment.ID != 1 || firstComment.PostID != 1 || firstComment.Author != "Comment author" || firstComment.Content != "Comment content" {
		t.Errorf("unexpected values in first comment: %+v", firstComment)
	}

	secondComment := comments[1]
	if secondComment.ID != 2 || secondComment.PostID != 1 || secondComment.Author != "Another author" || secondComment.Content != "Another content" {
		t.Errorf("unexpected values in second comment: %+v", secondComment)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProcessComments(t *testing.T) {
	commentMap := make(map[int]*model.Comment)
	commentMap[1] = &model.Comment{ID: 1, PostID: 1, Author: "Author 1", Content: "Content 1"}
	commentMap[2] = &model.Comment{ID: 2, PostID: 1, Author: "Author 2", Content: "Content 2", ParentID: &commentMap[1].ID}
	commentMap[3] = &model.Comment{ID: 3, PostID: 1, Author: "Author 3", Content: "Content 3"}

	comments := processComments(commentMap, 0, 10)

	if len(comments) != 2 {
		t.Errorf("expected length of comments list to be '2', got '%v'", len(comments))
	}

	if comments[0].ID != 1 || comments[1].ID != 3 {
		t.Errorf("unexpected comment IDs: got '%v' and '%v'", comments[0].ID, comments[1].ID)
	}

	if len(comments[0].Child) != 1 || comments[0].Child[0].ID != 2 {
		t.Errorf("unexpected child comment ID: got '%v'", comments[0].Child[0].ID)
	}

}

func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	ps := NewPostgresStore(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO posts").WithArgs("Test title", "Test content", "Test author", true).WillReturnRows(rows)

	post, err := ps.CreatePost("Test title", "Test content", "Test author")
	if err != nil {
		t.Errorf("error was not expected while creating post: %s", err)
	}

	if post.ID != 1 {
		t.Errorf("expected post ID to be '1', got '%d'", post.ID)
	}

	if post.Title != "Test title" {
		t.Errorf("expected post Title to be 'Test title', got '%s'", post.Title)
	}

	if post.Content != "Test content" {
		t.Errorf("expected post Content to be 'Test content', got '%s'", post.Content)
	}

	if post.Author != "Test author" {
		t.Errorf("expected post Author to be 'Test author', got '%s'", post.Author)
	}

	if post.CommentsEnabled != true {
		t.Errorf("expected post CommentsEnabled to be 'true', got '%v'", post.CommentsEnabled)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestCreateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ps := NewPostgresStore(db)

	mock.ExpectQuery("SELECT comments_enabled FROM posts WHERE id = \\$1").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"comments_enabled"}).AddRow(true))

	commentRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("INSERT INTO comments").WithArgs(1, "Comment author", "Comment content", nil).WillReturnRows(commentRows)

	comment, err := ps.CreateComment(1, "Comment author", "Comment content", nil)
	if err != nil {
		t.Errorf("error was not expected while creating comment: %s", err)
	}

	if comment == nil {
		t.Errorf("comment should not be nil")
		return
	}

	if comment.ID != 1 || comment.PostID != 1 || comment.Author != "Comment author" || comment.Content != "Comment content" {
		t.Errorf("unexpected values in comment: %+v", comment)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ps := NewPostgresStore(db)

	mock.ExpectExec("UPDATE posts SET title = \\$1, content = \\$2 WHERE id = \\$3").WithArgs("New title", "New content", 1).WillReturnResult(sqlmock.NewResult(1, 1))

	postRows := sqlmock.NewRows([]string{"id", "title", "content", "comments_enabled", "author"}).AddRow(1, "New title", "New content", true, "Test author")
	mock.ExpectQuery("^SELECT (.+) FROM posts WHERE id = \\$1").WithArgs(1).WillReturnRows(postRows)

	post, err := ps.UpdatePost(1, "New title", "New content")
	if err != nil {
		t.Errorf("error was not expected while updating post: %s", err)
	}

	if post == nil {
		t.Errorf("post should not be nil")
		return
	}

	if post.ID != 1 || post.Title != "New title" || post.Content != "New content" {
		t.Errorf("unexpected values in post: %+v", post)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ps := NewPostgresStore(db)

	mock.ExpectExec("UPDATE comments SET content = \\$1 WHERE id = \\$2").WithArgs("New content", 1).WillReturnResult(sqlmock.NewResult(1, 1))

	commentRows := sqlmock.NewRows([]string{"id", "post_id", "content", "author", "parent_id"}).AddRow(1, 1, "Test author", "New content", nil)
	mock.ExpectQuery("^SELECT (.+) FROM comments WHERE id = \\$1").WithArgs(1).WillReturnRows(commentRows)

	comment, err := ps.UpdateComment(1, "New content")
	if err != nil {
		t.Errorf("error was not expected while updating comment: %s", err)
	}

	if comment == nil {
		t.Errorf("comment should not be nil")
		return
	}

	if comment.ID != 1 || comment.Content != "New content" || comment.Author != "Test author" {
		t.Errorf("unexpected values in comment: %+v", comment)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDisableComments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ps := NewPostgresStore(db)

	mock.ExpectExec("UPDATE posts SET comments_enabled = false WHERE id = \\$1").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	err = ps.DisableComments(1)
	if err != nil {
		t.Errorf("error was not expected while disabling comments: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
