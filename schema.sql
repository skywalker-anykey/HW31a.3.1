/*
    Схема БД
*/

DROP TABLE IF EXISTS posts, authors;

CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    author_id INTEGER REFERENCES authors(id) NOT NULL DEFAULT 0,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
    published_at BIGINT NOT NULL DEFAULT extract(epoch from now())
);

-- Тестовые данные
INSERT INTO authors (id, name)
    VALUES (0, 'Аноним');
INSERT INTO posts (id, author_id, title, content, created_at,published_at)
VALUES
    (0, 0, 'Первая статья', 'Первый нах',0,0);

-- Тестовые данные без указания id и даты создание (установит текущую дату)
INSERT INTO authors (name)
VALUES
    ('Метеоролог');
INSERT INTO posts (author_id, title, content)
VALUES
    (1, 'Прогноз погоды на 28.09.2024','Солнечно. Температура: +50°C' );