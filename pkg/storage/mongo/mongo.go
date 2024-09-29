package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	db_posts_collection   = "posts"
	db_authors_collection = "authors"
	db_default_author     = "Anonymous"
)

// Store - хранилище данных.
type Store struct {
	db *mongo.Database
}

// New - конструктор объекта хранилища.
func New(connect, dbName string) (*Store, error) {
	mongoOpts := options.Client().ApplyURI(connect)
	db, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, err
	}
	return &Store{db: db.Database(dbName)}, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	return nil, nil
}

func (s *Store) AddPost(p storage.Post) error {
	// Выставляем метки времени
	t := time.Now().Unix()
	p.CreatedAt = t
	p.PublishedAt = t

	// Если не выбран автор, то указываем дефолтного
	if p.AuthorName == "" {
		p.AuthorName = db_default_author
	}

	// Получить ID автора
	authorId, err := s.GetAuthorId(p.AuthorName)
	if err != nil {
		return err
	}

	id, err := primitive.ObjectIDFromHex(authorId)
	if err != nil {
		panic(err)
	}

	newPost := struct {
		Post     storage.Post       `bson:"post"`
		AuthorID primitive.ObjectID `bson:"aid"`
	}{
		Post:     p,
		AuthorID: id,
	}

	log.Println(newPost)

	_, err = s.db.Collection(db_posts_collection).InsertOne(context.Background(), newPost)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdatePost(p storage.Post) error {
	return nil
}

func (s *Store) DeletePost(p storage.Post) error {
	return nil
}

// GetAuthorId - Найти автора по имени и вернуть его _id.ObjectID, если нет такого автора - создать и вернуть (используя CreateAuthorId).
func (s *Store) GetAuthorId(authorName string) (string, error) {
	authors := struct {
		Name     string             `bson:"name"`
		ObjectID primitive.ObjectID `bson:"_id"`
	}{}

	err := s.db.Collection(db_authors_collection).FindOne(context.Background(), bson.M{"name": authorName}).Decode(&authors)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			out, e := s.CreateAuthorId(authorName)
			if e != nil {
				return "", e
			}
			return out, nil
		} else {
			return "", err
		}
	}
	return authors.ObjectID.Hex(), nil
}

// CreateAuthorId - Создать нового автора и вернуть его _id.ObjectID
func (s *Store) CreateAuthorId(authorName string) (string, error) {
	// Добавляем автора
	dbAuthorName := bson.M{
		"name": authorName,
	}
	out, err := s.db.Collection(db_authors_collection).InsertOne(context.Background(), dbAuthorName)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(out.InsertedID), err
}
