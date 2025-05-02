# Pastebin

## Основные особенности проекта:
- Кодогенерация при помощи sqlc
- Хранение кода в Minio(S3)
- Кэширование популярных постов и ссылок на посты в Redis
- TTL на посты с автоочисткой из хранилища
- PostgreSQL с блокировками для конкурентного доступа к базе данных 

## Реализовано
- JWT-аутентификация и Auth-Middleware
- Логирование при помощи Logger-Middleware
- Кастомная обработка ошибок и возврат верных http кодом
- Валидация запросов на уровне бэкенда и на уровне БД


# Инструкция по запуску

Необходимо создать .env.{ENV} в папке configs. Значение ENV зависит от того, в каком окружении запускается приложения. Возможные окружения:
- prod
- dev
- test
Пример .env
```
MINIO_ROOT_USER=admin
MINIO_ROOT_PASSWORD=password123

POSTGRES_DB=pastebin_backend
POSTGRES_USER=admin
POSTGRES_PASSWORD=password123

PGADMIN_DEFAULT_EMAIL=admin@mail.ru
PGADMIN_DEFAULT_PASSWORD=password123
```
(не дописано)