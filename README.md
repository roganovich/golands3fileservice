# Goland s3 Fileservice API
Простое приложение на GO, предназначенное для работы с S3 сервером.

## Краткое описание ТЗ
Реализовать приложение позволяющее пользователям загружать и скачивать файлы с s3 сервера.

### Доступные операции
- Регистрация
- Авторизация JWT
- Работа с файлами
- Область видимости файла

### Роль пользователя
- Загрузка, скачивание своих файлов
- Шеринг файлов

### Роль Администратора
- Создание/Редактирование/Удаление карточки пользователя
- Создание/Редактирование/Удаление карточки файла

## Консольные команды

#### Получить все зависимости
```bash
go get .
```
#### Собрать приложение
```bash
go build -v .
```
#### Запустить приложение
```bash
go run .
```
#### Запустить контейнер Postgres
```bash
docker-compose --env-file .env.local up -d s3_postgres
```
#### Собрать контейнер API
```bash
docker compose build
```
#### Запустить контейнер API
```bash
docker compose --env-file .env.local up -d s3_app
```
#### Собрать контейнер API и запустить
```bash
docker-compose --env-file .env.local up --build
```

#### Зайти в контейнер
```bash
docker exec -it s3_app bash
```

#### Список контейнеров
```bash
docker ps -a
```
#### Список образов
```bash
docker images
```

### Миграции

#### Создать файл миграции
```bash
migrate create -ext sql -dir db/migration -seq add_role_to_user
```
#### Выполнить
```bash
migrate -path db/migration -database "postgresql://s3_postgres:s3_postgres@localhost:5432/s3_postgres?sslmode=disable" -verbose up
```

### Документация OpenAPI

#### Сгенирировать Swagger
```bash
swag init
```