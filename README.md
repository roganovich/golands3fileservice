# Goland s3 Fileservice API
Простое приложение на GO, предназначенное для работы с S3 сервером.

## Краткое описание ТЗ
Реализовать приложение позволяющее пользователям загружать и скачивать файлы с s3 сервера.

## S3 развернут в виде MinIO
API: http://172.25.0.3:9000  http://127.0.0.1:9000
WebUI: http://172.25.0.3:9001 http://127.0.0.1:9001
Docs: https://docs.min.io


### Доступные операции
- Регистрация
- Авторизация JWT
- Работа с файлами
- Область видимости файла

### Загрузка изображения
- После загрузки изображения мы должны его сжать, а оригинал сохранить
- При изменении изображении мы меняем файл из оригинала, а перезаписываем его копию

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

#### Собрать контейнер API
```bash
docker compose build
```

#### Запустить контейнер Postgres
```bash
docker compose --env-file .env.local up -d s3_postgres
```

#### Запустить контейнер S3
```bash
docker-compose --env-file .env.local up -d s3_minio
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