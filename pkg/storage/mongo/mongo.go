package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"errors"
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
	db           *mongo.Database
	nextPostID   int
	nextAuthorID int
}

// New - конструктор объекта хранилища.
func New(connect, dbName string) (*Store, error) {
	mongoOpts := options.Client().ApplyURI(connect)
	db, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, err
	}
	r := &Store{
		db:           db.Database(dbName),
		nextPostID:   0,
		nextAuthorID: 0,
	}

	err = r.SetIds()
	if err != nil {
		return nil, err
	}

	return r, nil
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
	p.AuthorID = authorId
	p.ID = s.nextPostID
	s.nextPostID++

	_, err = s.db.Collection(db_posts_collection).InsertOne(context.Background(), p)
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

// GetAuthorId - Найти автора по имени и вернуть его id, если нет такого автора - создать и вернуть (используя CreateAuthorId).
func (s *Store) GetAuthorId(authorName string) (int, error) {
	authors := struct {
		Name     string             `bson:"name"`
		ObjectID primitive.ObjectID `bson:"_id"`
		ID       int                `bson:"id"`
	}{}

	err := s.db.Collection(db_authors_collection).FindOne(context.Background(), bson.M{"name": authorName}).Decode(&authors)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			out, e := s.CreateAuthorId(authorName)
			if e != nil {
				return 0, e
			}
			return out, nil
		} else {
			return 0, err
		}
	}
	return authors.ID, nil
}

// CreateAuthorId - Создать нового автора и вернуть его id
func (s *Store) CreateAuthorId(authorName string) (int, error) {
	// Добавляем автора
	dbAuthorName := bson.M{
		"name": authorName,
		"id":   s.nextAuthorID,
	}
	_, err := s.db.Collection(db_authors_collection).InsertOne(context.Background(), dbAuthorName)
	if err != nil {
		return 0, err
	}
	out := s.nextAuthorID
	s.nextAuthorID++

	return out, err
}

// SetIds Устанавливает счетчики nextAuthorID и nextPostID
func (s *Store) SetIds() error {
	type Item struct {
		Id int
	}

	var results []Item

	// Получаем максимальное значение для id для колекции Авторов
	unsetStage := bson.D{{"$unset", bson.A{"_id", "name"}}}
	sortStage := bson.D{{"$sort", bson.D{{"id", -1}}}}
	limitStage := bson.D{{"$limit", 1}}

	cursor, err := s.db.Collection(db_authors_collection).Aggregate(context.TODO(), mongo.Pipeline{unsetStage, sortStage, limitStage})
	if err != nil {
		return err
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	var maxAuthID int

	for _, result := range results {
		maxAuthID = result.Id
	}
	s.nextAuthorID = maxAuthID + 1
	log.Println("SET nextAuthorID: ", s.nextAuthorID)

	// Получаем максимальное значение для id для колекции Постов
	unsetStage = bson.D{{"$unset", bson.A{"_id", "name"}}}
	sortStage = bson.D{{"$sort", bson.D{{"id", -1}}}}
	limitStage = bson.D{{"$limit", 1}}

	cursor, err = s.db.Collection(db_posts_collection).Aggregate(context.TODO(), mongo.Pipeline{unsetStage, sortStage, limitStage})
	if err != nil {
		log.Println("1 ", err)
		return err
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	var maxPostID int
	for _, result := range results {
		maxPostID = result.Id
	}
	s.nextPostID = maxPostID + 1
	log.Println("SET nextPostID: ", s.nextPostID)

	return nil
}
