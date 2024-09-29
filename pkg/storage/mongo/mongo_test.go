package mongo

import (
	"GoNews/pkg/storage"
	"log"
	"testing"
)

const (
	connect = "mongodb://192.168.1.20:27017/"
	dbName  = "news"
)

var S *Store

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
