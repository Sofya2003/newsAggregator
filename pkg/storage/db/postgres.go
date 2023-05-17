package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"skillfactory/36/pkg/storage"
)

type PostsDB struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, databaseUrl string) (*PostsDB, error) {
	db, err := pgxpool.Connect(ctx, databaseUrl)
	if err != nil {
		return nil, err
	}
	p := PostsDB{
		db: db,
	}
	return &p, nil
}

func (p *PostsDB) GetPosts(n int) ([]storage.Post, error) {
	rows, err := p.db.Query(context.Background(), `SELECT id, title, content, pubDate, link 
       FROM posts LIMIT $1;`, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var postsList []storage.Post
	for rows.Next() {
		var stPost storage.Post
		err = rows.Scan(
			&stPost.ID,
			&stPost.Title,
			&stPost.Content,
			&stPost.PubTime,
			&stPost.Link,
		)
		if err != nil {
			return nil, err
		}
		postsList = append(postsList, stPost)
	}
	return postsList, rows.Err()
}

func (p *PostsDB) AddPost(s storage.Post) error {
	err := p.db.QueryRow(context.Background(),
		`INSERT INTO posts (title, content, pubTime, link) VALUES 
             ($1, $2, $3, $4);`, s.Title, s.Content, s.PubTime, s.Link).Scan()
	return err
}
