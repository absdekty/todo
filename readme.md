# Todo REST API

REST API для управления задачами, реализованный на Go.

## Технологии

- Go 1.20+
- Chi (роутер)
- SQLite3 (БД)
- JWT (авторизация)
- Bcrypt (хэширование паролей)

## Особенности

- Clean Architecture
- Graceful Shutdown
- Метрики (total_requests, active_requests, errors)
- Middleware (логирование, recovery)
- Unit-тесты (model, service, repository)

## Эндпоинты

### Публичные

| Метод |   Эндпоинт  | Описание |
|-------|-------------|----------|
| POST  | `/register` | Регистрация пользователя |
| POST  | `/login`    | Аутентификация, получение JWT токена |
| POST  | `/refresh`  | Получение нового access-токена |
| POST  | `/logout`   | Выход с сессии, очистка куки |
| GET   | `/metrics`  | Метрики сервера |

### Защищенные (требуют Bearer токен)

| Метод  |    Эндпоинт   | Описание |
|--------|---------------|----------|
| GET    | `/tasks`      | Получить все задачи пользователя |
| GET    | `/tasks/{id}` | Получить конкретную задачу |
| POST   | `/tasks`      | Создать задачу |
| PUT    | `/tasks/{id}` | Полностью обновить задачу |
| PATCH  | `/tasks/{id}` | Частично обновить задачу |
| DELETE | `/tasks/{id}` | Удалить задачу |

## Пример использования (Использовать в одном bash-терминале)

#### Регистрация
```
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name":"testuser","password":"password123"}'
```

#### Аутентификация+получение токена
```
TOKEN=$(curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"name":"testuser","password":"password123"}' \
  -c cookies.txt \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
```

#### Получение access-токена по истечению
```
TOKEN=$(curl -X POST http://localhost:8080/refresh \
  -b cookies.txt \
  -c cookies.txt \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
```

#### Выход с сессии
```
curl -X POST http://localhost:8080/logout \
  -c cookies.txt \
  -b cookies.txt
```

#### Создание задачи
```
curl -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy groceries","desc":"Milk, eggs, bread"}'
```

#### Получить все задачи
```
curl -X GET http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN"
```

#### Получить конкретную задачу по ID
```
curl -X GET http://localhost:8080/ID-ЗАДАЧИ \
  -H "Authorization: Bearer $TOKEN"
```

#### Обновить всю задачу
```
curl -X PUT http://localhost:8080/tasks/ID-ЗАДАЧИ \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy groceries","desc":"Milk, eggs, bread, cheese","completed":true}'
```

#### Обновить любые поля(все виды полей в запросе PUT)
```
curl -X PATCH http://localhost:8080/tasks/ID-ЗАДАЧИ \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"completed":true}'
```

#### Удалить задачу
```
curl -X DELETE http://localhost:8080/tasks/ID-ЗАДАЧИ \
  -H "Authorization: Bearer $TOKEN"
```