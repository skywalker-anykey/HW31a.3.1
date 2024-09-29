package postgres

import (
	"GoNews/pkg/storage"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Store - хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

// New - конструктор объекта хранилища.
func New(connect string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), connect)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
SELECT
		posts.id AS id,
		title,
		content,
		author_id,
		authors.name,
		created_at,
		published_at
	FROM posts
	    JOIN authors ON author_id=authors.id
	ORDER BY id;
`)
	if err != nil {
		return nil, err
	}

	var posts []storage.Post
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var post storage.Post
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.AuthorName,
			&post.CreatedAt,
			&post.PublishedAt,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		posts = append(posts, post)
	}
	// ВАЖНО не забыть проверить rows.Err()
	return posts, rows.Err()
}

func (s *Store) AddPost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
	INSERT INTO posts
		(title, content, author_id)	
	VALUES 
		($1,$2,$3);
`,
		&p.Title,
		&p.Content,
		&p.AuthorID,
	)
	return err
}

func (s *Store) UpdatePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
	UPDATE posts
	SET
		title = $2,
		content = $3,
		author_id = $4,
		created_at =$5,
		published_at = $6
	WHERE id = $1;
`,
		&p.ID,
		&p.Title,
		&p.Content,
		&p.AuthorID,
		&p.CreatedAt,
		&p.PublishedAt,
	)
	return err
}

func (s *Store) DeletePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
	DELETE FROM posts WHERE id = $1;`, &p.ID,
	)
	return err
}
