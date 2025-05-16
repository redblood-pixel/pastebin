# Pastebin

## Основные особенности проекта:
- Хранение кода в Minio(S3)
- "Ленивое" удаление постов из БД и S3
- PostgreSQL с блокировками для конкурентного доступа к базе данных
- Unit и интеграционное тестирование проекта

## Реализовано
- JWT-аутентификация и Auth-Middleware
- Логирование при помощи Logger-Middleware
- Кастомная обработка ошибок и возврат верных http кодом
- Валидация запросов на уровне бэкенда и на уровне БД
- Возможность создавать TTL на посты
- Возможность создавать пароли на посты
- Возможность задать для поста режим "удалить после прочтения"
- Получение постов пользователя с фильтрацией, сортировкой и пагинацией


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

После чего необходимо создать yaml конфиг с настройками бэкенда

Пример настроек
```
env: dev

jwt:
  access_token_ttl: 60m
  refresh_token_ttl: 168h
  signing_key: 123

http:
  port: 8080
  write_timeout: 10s
  read_timeout: 10s

postgres:
  username: admin
  password: password
  host: localhost
  port: 5436
  db_name: backend
  ssl_mode: disable

minio:
  bucket_name: pastes
  access_key: Pastebin_backend
  secret_key: password1
  use_secure: false
```
(не дописано)