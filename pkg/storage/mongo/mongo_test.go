package mongo

import (
	"GoNews/pkg/storage"
	"fmt"
	"log"
	"testing"
)

const (
	connect = "mongodb://192.168.1.20:27017/"
	dbName  = "news"
)

//var S *Store

func TestNew(t *testing.T) {
	type args struct {
		connect string
		dbName  string
	}
	tests := []struct {
		name    string
		args    args
		want    *Store
		wantErr bool
	}{
		{
			name: "Test New",
			args: args{
				connect: connect,
				dbName:  dbName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.connect, tt.args.dbName)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
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
			name: "Test AddPost_1",
			args: args{
				p: storage.Post{
					Title:       "Add Post",
					Content:     "Add Content",
					AuthorName:  "Add Author",
					CreatedAt:   0,
					PublishedAt: 0,
				},
			},
			wantErr: false,
		},
		{
			name: "Test AddPostdb_default_author",
			args: args{
				p: storage.Post{
					Title:       "Add Post",
					Content:     "Add Content",
					CreatedAt:   0,
					PublishedAt: 0,
				},
			},
			wantErr: false,
		},
	}

	s, err := New(connect, dbName)
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

func TestStore_DeletePost(t *testing.T) {
	s, err := New(connect, dbName)
	if err != nil {
		log.Fatal(err)
	}
	lastPostID := s.nextPostID - 1

	type args struct {
		p storage.Post
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete Last Post",
			args: args{
				p: storage.Post{
					ID: lastPostID,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.DeletePost(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("DeletePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_Posts(t *testing.T) {
	s, err := New(connect, dbName)
	if err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Get All Posts",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Posts()
			if (err != nil) != tt.wantErr {
				t.Errorf("Posts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}

func TestStore_UpdatePost(t *testing.T) {
	s, err := New(connect, dbName)
	if err != nil {
		log.Fatal(err)
	}
	lastPostID := s.nextPostID - 1

	type args struct {
		p storage.Post
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Update Last Post",
			args: args{
				p: storage.Post{
					ID:      lastPostID,
					Title:   "[update] Title",
					Content: "[update] Content",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.UpdatePost(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
