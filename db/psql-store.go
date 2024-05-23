package db

import (
	"PostCommentService/graph/model"
	"database/sql"
	"errors"
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

func (s *PostgresStore) GetPost(id, offset, limit int) (*model.Post, error) {
	row := s.db.QueryRow("SELECT id, title, content, comments_enabled, author FROM posts WHERE id = $1", id)

	var p model.Post
	if err := row.Scan(&p.ID, &p.Title, &p.Content, &p.CommentsEnabled, &p.Author); err != nil {
		return nil, err
	}

	comments, err := s.GetComments(p.ID, offset, limit)
	if err != nil {
		return nil, err
	}

	p.Comments = comments

	return &p, nil
}

func (s *PostgresStore) GetComments(postID, offset, limit int) ([]*model.Comment, error) {
	rows, err := s.db.Query("SELECT id, post_id, content, author, parent_id FROM comments WHERE post_id = $1 ORDER BY id ASC ", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	commentMap := make(map[int]*model.Comment)

	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.Author, &c.Content, &c.ParentID); err != nil {
			return nil, err
		}
		commentMap[c.ID] = &c
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return processComments(commentMap, offset, limit), nil
}

func processComments(commentMap map[int]*model.Comment, offset, limit int) []*model.Comment {
	// Обработка комментариев
	for _, comment := range commentMap {
		if comment.ParentID != nil {
			parentComment, ok := commentMap[*comment.ParentID]
			if ok {
				parentComment.Child = append(parentComment.Child, comment)
			}
		}
	}

	// Возвращаем только комментарии верхнего уровня с учетом пагинации
	var topLevelComments []*model.Comment
	for _, comment := range commentMap {
		if comment.ParentID == nil {
			topLevelComments = append(topLevelComments, comment)
		}
	}

	// Применяем пагинацию
	if len(topLevelComments) <= offset {
		return []*model.Comment{}
	}

	end := offset + limit
	if end > len(topLevelComments) {
		end = len(topLevelComments)
	}

	return topLevelComments[offset:end]
}

func (s *PostgresStore) GetComment(id int) (*model.Comment, error) {
	row := s.db.QueryRow("SELECT id, post_id, author, content, parent_id FROM comments WHERE id = $1", id)

	var c model.Comment
	if err := row.Scan(&c.ID, &c.PostID, &c.Author, &c.Content, &c.ParentID); err != nil {
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
	if len(content) > 2000 {
		return nil, errors.New("comment is too long")
	}

	// Проверяем, разрешены ли комментарии для этого поста
	var commentsEnabled bool
	err := s.db.QueryRow("SELECT comments_enabled FROM posts WHERE id = $1", postID).Scan(&commentsEnabled)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	if !commentsEnabled {
		return nil, errors.New("comments are disabled for this post")
	}

	var c model.Comment
	err = s.db.QueryRow("INSERT INTO comments(post_id, author, content, parent_id) VALUES($1, $2, $3, $4) RETURNING id", postID, author, content, parentID).Scan(&c.ID)
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

	return s.GetPost(id, 0, 0)
}

func (s *PostgresStore) UpdateComment(id int, content string) (*model.Comment, error) {
	if len(content) > 2000 {
		return nil, errors.New("comment is too long")
	}

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
