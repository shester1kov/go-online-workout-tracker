# Проект: Онлайн-трекер тренировок

## Описание

Этот проект представляет собой серверную часть онлайн-трекера тренировок. API разработано на Go и предоставляет функционал для управления пользователями, упражнениями, категориями и ролями. Проект использует авторизацию через JWT, кеширование в Redis и взаимодействует с PostgreSQL.

## Технологии

- Go (Chi)
- PostgreSQL
- Redis
- Docker
- Swagger
- JWT
- Nginx
- Fail2Ban (защита от брутфорса)

## Структура проекта

```
├── cmd/
│   └── app/
│       └── main.go             # Точка входа
│
├── docs/
│   ├── docs.go                 # Настройки Swagger
│   ├── swagger.json            # Сгенерированный Swagger JSON
│   ├── swagger.yaml            # Сгенерированный Swagger YAML
│
├── internal/
│   ├── apperrors/              # Глобальная обработка ошибок
│   ├── appmiddlewares/         # Middleware (аутентификация, роли)
│   ├── auth/                   # Логика JWT
│   ├── config/                 # Конфигурация проекта
│   ├── db/                     # Подключение к БД и миграции
│   ├── handlers/               # Обработчики API
│   ├── models/                 # Модели данных
│   ├── repository/             # Логика работы с БД
│   ├── server/                 # Настройки сервера и маршрутов
│   ├── services/               # Бизнес-логика
│   └── utils/                  # Утилиты
│
├── migrations/                 # SQL-миграции
├── go.mod
├── go.sum
└── Dockerfile
```

## Запуск проекта

### Требования

- Установленный Docker и Docker Compose

### Шаги

1. Склонировать репозиторий:
```sh
   git clone <URL репозитория>
   ```
2. Собрать и запустить контейнеры:

   ```sh
   docker-compose up
   ```

## Безопасность

- Авторизация с использованием JWT
- Ограничение доступа по ролям (admin, user, moderator, trainer)
- Кеширование данных в Redis для оптимизации запросов
- Защита от брутфорс-атак с помощью Fail2Ban

## Тестирование

- Проверка работы API через Swagger (доступен по `/api/v1/swagger/index.html`)

## Логи

- Логирование запросов и ошибок
- Доступны через stdout или файлы логов при запуске через Docker

## Выводы

Проект реализует основные принципы чистой архитектуры, использует минимальные зависимости и обеспечивает удобное взаимодействие с базой данных. Реализована система ролей, кеширование и защита API.
