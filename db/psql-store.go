package db

import (
	"PostCommentService/graph/model"
	"database/sql"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{
		db: db,
	}
}

func (s *PostgresStore) GetPosts() ([]*model.Post, error) {
	rows, err := s.db.Query("SELECT id, title, content, comments_enabled, author FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CommentsEnabled, &p.Author); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostgresStore) GetPost(id int) (*model.Post, error) {
	row := s.db.QueryRow("SELECT id, title, content, comments_enabled, author FROM posts WHERE id = $1", id)

	var p model.Post
	if err := row.Scan(&p.ID, &p.Title, &p.Content, &p.CommentsEnabled, &p.Author); err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *PostgresStore) GetComments() ([]*model.Comment, error) {
	rows, err := s.db.Query("SELECT id, post_id, author, content, parent_id, child_ids FROM comments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.Author, &c.Content, &c.ParentID, &c.Child); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *PostgresStore) GetComment(id int) (*model.Comment, error) {
	row := s.db.QueryRow("SELECT id, post_id, author, content, parent_id, child_ids FROM comments WHERE id = $1", id)

	var c model.Comment
	if err := row.Scan(&c.ID, &c.PostID, &c.Author, &c.Content, &c.ParentID, &c.Child); err != nil {
		return nil, err
	}

	return &c, nil
}

func (s *PostgresStore) CreatePost(title, content, author string) (*model.Post, error) {
	var p model.Post
	err := s.db.QueryRow("INSERT INTO posts(title, content, author,comments_enabled) VALUES($1, $2, $3, $4) RETURNING id",
		title, content, author, true).Scan(&p.ID)
	if err != nil {
		return nil, err
	}

	p.Title = title
	p.Content = content
	p.Author = author
	p.CommentsEnabled = true

	return &p, nil
}

func (s *PostgresStore) CreateComment(postID int, author, content string, parentID *int) (*model.Comment, error) {
	var c model.Comment
	err := s.db.QueryRow("INSERT INTO comments(post_id, author, content, parent_id) VALUES($1, $2, $3, $4) RETURNING id", postID, author, content, parentID).Scan(&c.ID)
	if err != nil {
		return nil, err
	}

	c.PostID = postID
	c.Author = author
	c.Content = content
	c.ParentID = parentID

	return &c, nil
}

func (s *PostgresStore) UpdatePost(id int, title, content string) (*model.Post, error) {
	_, err := s.db.Exec("UPDATE posts SET title = $1, content = $2 WHERE id = $3", title, content, id)
	if err != nil {
		return nil, err
	}

	return s.GetPost(id)
}

func (s *PostgresStore) UpdateComment(id int, content string) (*model.Comment, error) {
	_, err := s.db.Exec("UPDATE comments SET content = $1 WHERE id = $2", content, id)
	if err != nil {
		return nil, err
	}

	return s.GetComment(id)
}

func (s *PostgresStore) DisableComments(postID int) error {
	_, err := s.db.Exec("UPDATE posts SET comments_enabled = false WHERE id = $1", postID)
	return err
}