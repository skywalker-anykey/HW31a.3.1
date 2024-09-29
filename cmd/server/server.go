package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/memdb"
	"GoNews/pkg/storage/postgres"
	"log"
	"net/http"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server

	// Создаём объекты баз данных.
	//
	// БД в памяти.
	db1 := memdb.New()

	// Реляционная БД Postgres SQL.
	db2, err := postgres.New("postgres://sandbox:sandbox@localhost:5432/news")
	if err != nil {
		log.Fatal(err)
	}
	// Документная БД MongoDB.
	//db3, err := mongo.New("mongodb://server.domain:27017/")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//_, _, _ = db1, db2, db3
	_, _ = db1, db2

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db2

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	_ = http.ListenAndServe(":8080", srv.api.Router())
}
