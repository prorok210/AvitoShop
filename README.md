# AvitoShop

API сервис для работы с товарами, балансом пользователей и транзакциями.

## Технологии

- Go 1.22

- Echo framework

- PostgreSQL 13

- Docker & Docker Compose

- Swagger (OpenAPI 3.0)

## Запуск приложения

### Production окружение

Сборка:

```bash
docker-compose  --profile  prod  build
```

Запуск:

```bash
docker-compose  --profile  prod  up  -d
```

### Тестовое окружение

Сборка и запуск e2e тестов:

```bash
docker-compose  --profile  test  build
docker-compose  --profile  test  up  --abort-on-container-exit  --exit-code-from  e2e-tests
```

Юнит-тесты запускаются автоматически при сборке проекта

Текущее покрытие тестами:

```bash
coverage:  53.7%  of  statements  in  ./...
```

### API Документация

Swagger UI доступен по адресу: http://127.0.0.1:8080/swagger/

**Важно**: Для корректного отображения документации выберите схему: http://localhost:8080/docs/swagger.json

Для генерации/обновления Swagger документации:

```bash
swag  init  -g  cmd/avito_shop_service/main.go
```

### Основной функционал

- Регистрация и аутентификация пользователей

- Просмотр баланса и истории операций

- Покупка товаров из каталога

- Перевод монет между пользователями

- E2E тестирование основных сценариев
