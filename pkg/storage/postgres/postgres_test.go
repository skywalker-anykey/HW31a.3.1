package postgres

import (
	"GoNews/pkg/storage"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"testing"
)

var s *Store

func TestStore_Posts(t *testing.T) {
	s, err := New("postgres://sandbox:sandbox@localhost:5432/news")
	if err != nil {
		log.Fatal(err)
	}

	data, err := s.Posts()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestStore_UpdatePost(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		p storage.Post
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "UpdatePostID0",
			args: args{
				p: storage.Post{
					ID:          0,
					Title:       "[Обновлено] Первая статья",
					Content:     "[Обновлено] Первый нах",
					AuthorID:    1,
					AuthorName:  "",
					CreatedAt:   1727594802,
					PublishedAt: 1727594802,
				},
			},
			wantErr: false,
		},
	}
	s, err := New("postgres://sandbox:sandbox@localhost:5432/news")
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.UpdatePost(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_DeletePost(t *testing.T) {
	type args struct {
		p storage.Post
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "DeletePostID1",
			args: args{
				p: storage.Post{
					ID: 1,
				},
			},
			wantErr: false,
		},
	}
	s, err := New("postgres://sandbox:sandbox@localhost:5432/news")
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.DeletePost(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("DeletePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_AddPost(t *testing.T) {
	type args struct {
		p storage.Post
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "AddPost 1",
			args: args{
				p: storage.Post{
					Title:   "[Новый] Ура мы вернулись",
					Content: "Снова пишем статьи анонимно",
				},
			},
			wantErr: false,
		},
		{
			name: "AddPost 2",
			args: args{
				p: storage.Post{
					Title:    "[Новый] Прогноз погоды на завтра",
					Content:  "Дождь. Температура: -50°C",
					AuthorID: 1,
				},
			},
			wantErr: false,
		},
	}

	s, err := New("postgres://sandbox:sandbox@localhost:5432/news")
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.AddPost(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("AddPost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
