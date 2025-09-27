# GoIda - Система управления статьями

Веб-приложение для управления статьями с авторизацией пользователей, построенное на Go (бэкенд) и Vue.js (фронтенд).

## Описание проекта

GoIda - это система управления статьями, которая позволяет пользователям создавать, редактировать и удалять статьи. Система включает в себя авторизацию пользователей, управление ролями и административные функции.

## Возможности

- **Авторизация и регистрация** - вход в систему с логином и паролем, регистрация новых пользователей
- **Управление статьями** - создание, просмотр, редактирование и удаление статей
- **Управление пользователями** - просмотр списка пользователей (только для админов)
- **Логи запросов** - отслеживание всех API запросов в реальном времени

## Быстрый старт

### Предварительные требования

- Docker и Docker Compose

### Запуск

1. **Клонируйте репозиторий:**
   ```bash
   git clone <repository-url>
   cd GoIda
   ```

2. **Запустите все сервисы:**
   ```bash
   docker-compose up -d
   ```

3. **Откройте браузер:**
   - Фронтенд: http://localhost:3000
   - API: http://localhost:8080

## Тестовые аккаунты

- **Админ**: логин `admin`, пароль `password`
- **Пользователь**: логин `user`, пароль `password`

## Функции

### Авторизация
- **Вход**: Введите логин и пароль, нажмите "Войти"
- **Регистрация**: Переключитесь на вкладку "Регистрация", заполните форму
- **Выход**: Нажмите "Выйти" для завершения сессии

### Статьи
- **Создание**: Заполните заголовок и содержимое, нажмите "Создать статью"
- **Просмотр**: Нажмите "Обновить список" для загрузки статей
- **Редактирование**: Доступно только для своих статей (или для админов)
- **Удаление**: Доступно только для своих статей (или для админов)

### Пользователи (только для админов)
- Нажмите "Загрузить пользователей" для просмотра списка

### Логи
- Все API запросы отображаются в разделе "Логи запросов"

## API Endpoints

### Публичные запросы

#### Авторизация

**POST** `/api/auth/login` - авторизация пользователя

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Parameters: `{"login":"admin","password":"password"}` | **Success:** *Пользователь найден*<br/>Status: 200/OK<br/>Content-type: application/json<br/>Body: `{"token":"jwt_token","user":{"id":1,"email":"admin@example.com","name":"Admin","role":{"name":"admin"}}}`<br/>**Denied:** *Неверные данные*<br/>Status: 401 |

#### Регистрация пользователя

**POST** `/api/users` - регистрация нового пользователя

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Parameters: `{"name":"Имя","email":"email@example.com","login":"username","password":"password"}` | **Success:** *Пользователь создан*<br/>Status: 201/Created<br/>Content-type: application/json<br/>Body: `{"message":"User created successfully","user":{"id":1,"email":"email@example.com","name":"Имя","role":{"name":"user"}},"login":"username"}`<br/>**Denied:** *Логин уже занят*<br/>Status: 409<br/>**Validation Error:** *Неверные данные*<br/>Status: 422 |

#### Список статей

**GET** `/api/articles` - получение списка статей

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Query parameters: `?limit=10&offset=0` | **Success:** *Статьи найдены*<br/>Status: 200/OK<br/>Content-type: application/json<br/>Body: `[{"id":1,"title":"Заголовок","content":"Содержимое","author_id":1,"author_name":"Автор","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"}]` |

#### Информация о статье

**GET** `/api/articles/{id}` - получение информации о статье

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Parameters: id статьи в URL | **Success:** *Статья найдена*<br/>Status: 200/OK<br/>Content-type: application/json<br/>Body: `{"id":1,"title":"Заголовок","content":"Содержимое","author_id":1,"author_name":"Автор","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"}`<br/>**Denied:** *Статья не найдена*<br/>Status: 404 |

#### Статьи пользователя

**GET** `/api/users/{authorId}/articles` - получение статей конкретного пользователя

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Parameters: authorId в URL<br/>Query parameters: `?limit=10&offset=0` | **Success:** *Статьи найдены*<br/>Status: 200/OK<br/>Content-type: application/json<br/>Body: `[{"id":1,"title":"Заголовок","content":"Содержимое","author_id":1,"author_name":"Автор","created_at":"2024-01-01T00:00:00Z"}]` |

### Авторизованные запросы

Для выполнения авторизованных запросов необходимо передать Bearer токен в заголовке Authorization.

#### Профиль пользователя

**GET** `/api/auth/profile` - получение профиля текущего пользователя

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Authorization: Bearer <токен> | **Success:** *Профиль получен*<br/>Status: 200/OK<br/>Content-type: application/json<br/>Body: `{"id":1,"email":"email@example.com","name":"Имя","role":{"name":"user"}}`<br/>**Denied:** *Неверный токен*<br/>Status: 401 |

#### Создание статьи

**POST** `/api/articles` - создание новой статьи

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Authorization: Bearer <токен><br/>Parameters: `{"title":"Заголовок","content":"Содержимое статьи"}` | **Success:** *Статья создана*<br/>Status: 201/Created<br/>Content-type: application/json<br/>Body: `{"id":1,"title":"Заголовок","content":"Содержимое","author_id":1,"author_name":"Автор","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"}`<br/>**Validation Error:** *Неверные данные*<br/>Status: 422 |

#### Редактирование статьи

**PUT** `/api/articles/{id}` - редактирование статьи

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Authorization: Bearer <токен><br/>Parameters: `{"title":"Новый заголовок","content":"Новое содержимое"}` | **Success:** *Статья обновлена*<br/>Status: 200/OK<br/>Content-type: application/json<br/>Body: `{"id":1,"title":"Новый заголовок","content":"Новое содержимое","author_id":1,"author_name":"Автор","updated_at":"2024-01-01T00:00:00Z"}`<br/>**Denied:** *Нет прав*<br/>Status: 403<br/>**Not Found:** *Статья не найдена*<br/>Status: 404 |

#### Удаление статьи

**DELETE** `/api/articles/{id}` - удаление статьи

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Authorization: Bearer <токен><br/>Parameters: id статьи в URL | **Success:** *Статья удалена*<br/>Status: 204/No Content<br/>**Denied:** *Нет прав*<br/>Status: 403<br/>**Not Found:** *Статья не найдена*<br/>Status: 404 |

#### Получение пользователя

**GET** `/api/users/{id}` - получение информации о пользователе

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Authorization: Bearer <токен><br/>Parameters: id пользователя в URL | **Success:** *Пользователь найден*<br/>Status: 200/OK<br/>Content-type: application/json<br/>Body: `{"id":1,"email":"email@example.com","name":"Имя","role":{"name":"user"}}`<br/>**Not Found:** *Пользователь не найден*<br/>Status: 404 |

### Административные запросы

Административные запросы доступны только пользователям с ролью "admin".

#### Список всех пользователей

**GET** `/api/admin/users` - получение списка всех пользователей

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Authorization: Bearer <токен><br/>Query parameters: `?limit=10&offset=0` | **Success:** *Пользователи найдены*<br/>Status: 200/OK<br/>Content-type: application/json<br/>Body: `[{"id":1,"email":"email@example.com","name":"Имя","role_id":1,"role":{"id":1,"name":"user","description":"Обычный пользователь"},"created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"}]`<br/>**Denied:** *Нет прав*<br/>Status: 403 |

#### Список ролей

**GET** `/api/admin/roles` - получение списка ролей

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Authorization: Bearer <токен> | **Success:** *Роли найдены*<br/>Status: 200/OK<br/>Content-type: application/json<br/>Body: `[{"id":1,"name":"user","description":"Обычный пользователь","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"},{"id":2,"name":"admin","description":"Администратор","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"}]` |

#### Информация о роли

**GET** `/api/admin/roles/{id}` - получение информации о роли

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Authorization: Bearer <токен><br/>Parameters: id роли в URL | **Success:** *Роль найдена*<br/>Status: 200/OK<br/>Content-type: application/json<br/>Body: `{"id":1,"name":"user","description":"Обычный пользователь","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"}`<br/>**Not Found:** *Роль не найдена*<br/>Status: 404 |

### Дополнительные эндпоинты

#### Управление учетными данными

**POST** `/api/auth/credentials` - создание учетных данных для пользователя

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Parameters: `{"user_id":1,"login":"username","password":"password"}` | **Success:** *Учетные данные созданы*<br/>Status: 201/Created<br/>Content-type: application/json<br/>Body: `{"id":1,"user_id":1,"login":"username","created_at":"2024-01-01T00:00:00Z"}` |

**GET** `/api/auth/credentials/{userId}` - получение учетных данных пользователя

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Parameters: userId в URL | **Success:** *Учетные данные найдены*<br/>Status: 200/OK<br/>Content-type: application/json<br/>Body: `{"id":1,"user_id":1,"login":"username","created_at":"2024-01-01T00:00:00Z"}` |

**PUT** `/api/auth/credentials/{userId}` - обновление учетных данных пользователя

| Request | Response |
| :---- | :---- |
| Content-type: application/json<br/>Parameters: `{"login":"new_username","password":"new_password"}` | **Success:** *Учетные данные обновлены*<br/>Status: 200/OK<br/>Content-type: application/json<br/>Body: `{"id":1,"user_id":1,"login":"new_username","updated_at":"2024-01-01T00:00:00Z"}` |

## Роли пользователей

- **user** - обычный пользователь (может создавать и редактировать только свои статьи)
- **admin** - администратор (может управлять всеми статьями и пользователями)

## Валидация данных

### Статьи
- **title** - обязательное поле, минимум 3 символа
- **content** - обязательное поле, минимум 10 символов

### Пользователи
- **name** - обязательное поле, минимум 2 символа
- **email** - обязательное поле, валидный email
- **login** - обязательное поле, минимум 3 символа
- **password** - обязательное поле, минимум 6 символов

### Учетные данные
- **login** - обязательное поле, минимум 3 символа
- **password** - обязательное поле

## Коды ошибок

- **200** - OK (успешный запрос)
- **201** - Created (ресурс создан)
- **204** - No Content (ресурс удален)
- **400** - Bad Request (неверный запрос)
- **401** - Unauthorized (не авторизован)
- **403** - Forbidden (нет прав доступа)
- **404** - Not Found (ресурс не найден)
- **409** - Conflict (конфликт, например, логин уже занят)
- **422** - Unprocessable Entity (ошибка валидации)
- **500** - Internal Server Error (внутренняя ошибка сервера)

## Устранение неполадок

### CORS ошибки
Убедитесь, что бэкенд запущен и CORS настроен правильно.

### API недоступен
Проверьте статус контейнеров:
```bash
docker ps
```

### База данных
Проверьте логи миграций:
```bash
docker logs goida-liquibase
```

### Порт занят
Если порт 3000 или 8080 занят, измените порты в `docker-compose.yml`.

## Разработка

### Локальная разработка

1. **Бэкенд:**
   ```bash
   go run main.go
   ```

2. **Фронтенд:**
   ```bash
   cd goida-frontend
   python -m http.server 3000
   ```

### Пересборка контейнеров

```bash
# Пересборка всех сервисов
docker-compose build --no-cache

# Пересборка конкретного сервиса
docker-compose build --no-cache app
docker-compose build --no-cache frontend
```

### Логи

```bash
# Все сервисы
docker-compose logs

# Конкретный сервис
docker-compose logs app
docker-compose logs frontend
```