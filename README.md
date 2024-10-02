# HW31a.3.1 (go-db-practice-gonews)

## Задание:

- [X] Разработать схему БД PostgreSQL в форме SQL-запроса. Запрос должен быть помещён в файл schema.sql в корневой каталог проекта.

- [X] По аналогии с пакетом "memdb" разработать пакет "postgres" для поддержки базы данных под управлением СУБД PostgreSQL.

- [X] По аналогии с пакетом "memdb" разработать пакет "mongo" для поддержки базы данных под управлением СУБД MongoDB.

## Замечания от ментора:

- [ ] В функции New вы подключаетесь к базе данных через mongo.Connect, но не проверяете успешность подключения с помощью метода Ping. После подключения нужно убедиться, что соединение действительно установлено, вызвав метод Ping на объекте клиента

- [ ] Переменные nextPostID и nextAuthorID хранятся в памяти приложения. Это значит, что при перезапуске сервера или приложения эти значения будут сбрасываться. Возможно, стоило бы использовать автоматическое инкрементирование ObjectID, встроенное в MongoDB, или хранить текущие значения идентификаторов в самой базе данных.

- [X] В функции Posts при возникновении ошибки вы выводите ее через log.Println, но не возвращаете ошибку пользователю. Если произошла ошибка в Find, то лучше вернуть ее в ответ
