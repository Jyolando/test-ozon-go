# Сервис по созданию сокращённых ссылок

---

## F.A.Q

### Что это?

Тестовое задание для Ozon Fintech представляющее из себя сервис для создания сокрщённых ссылок.

### Что используется?

- Golang
- gRPC
- PostgreSQL
- Docker

### В чём суть задания?

Реализовать упакованное в Docker-образ решение использующее gRPC, два метода хранения данных (in-memory и postgreSQL - выбирается во время запуска сервиса) и покрытый тестами функционал.

---

## Как взаимодействовать?

1. Добавить .env файл в корень проекта со следующими аргументами:
   * POSTGRES_HOST
   * POSTGRES_PORT
   * POSTGRES_DB
   * POSTGRES_USER
   * POSTGRES_PASSWORD


2. Отправить запросы:

    **Создание сокращённой ссылки:**
    ```protobuf
   AddURL{url: string}
   ```
   
   **Возможные ответы:**
   ```json lines
   0 - OK
   {
      "url": {
        "originalURL": "https://example.com",
        "shortenedURL": "LzxAsRmQx3"
      }
   }
   
   3 - INVALID LINK
   
   13 - SERVER ERROR
   ```
   
   **Получение оригинальной ссылки из сокращённой:**
   ```protobuf
   GetURL{url: string}
   ```

   **Возможные ответы:**
   ```json lines
   0 - OK
   {
      "url": {
        "originalURL": "https://example.com",
        "shortenedURL": "LzxAsRmQx3"
      }
   }
   
   3 - INVALID LINK
   
   5 - ELEMENT NOT FOUND
   
   13 - SERVER ERROR
   ```


## Как запустить?

```shell
git clone https://github.com/jyolando/test-ozon-go.git && cd test-ozon-go

# Запуск тестов:
make test

# Запуск с in-memory хранилищем:
make memory

# Запуск с хранилищем Postgres:
make postgresql
```
