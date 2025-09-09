# GoIda

## Запуск
```bash
cp env.example .env
docker compose up -d

curl http://localhost:8080/users
```

## API
- `GET /users` - список
- `POST /users` - создать
- `GET /users/{id}` - получить
- `PUT /users/{id}` - обновить
- `DELETE /users/{id}` - удалить
