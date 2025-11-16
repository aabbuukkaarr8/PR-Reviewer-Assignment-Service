# PR Reviewer Assignment Service

Микросервис для автоматического назначения ревьюверов на Pull Request'ы, управления командами и участниками.

##  Описание

Сервис предоставляет HTTP API для:
- Создания и управления командами разработчиков
- Автоматического назначения ревьюверов на PR (до 2 активных ревьюверов из команды автора)
- Переназначения ревьюверов
- Управления активностью пользователей
- Получения списка PR, назначенных конкретному пользователю

##  Быстрый старт

### Требования

- Go 1.25+ (для локальной разработки)
- Docker и Docker Compose (для запуска через Docker)
- PostgreSQL 16+ (если запускаете БД отдельно)

### Запуск через Docker (рекомендуется)

1. **Клонируйте репозиторий:**
   ```bash
   git clone <repository-url>
   cd PR-Reviewer-Assignment-Service
   ```

2. **Запустите сервисы:**
   ```bash
   docker-compose up -d
   ```

   Это автоматически:
   - Соберет образы для БД и API
   - Запустит PostgreSQL с применением миграций
   - Запустит API сервис на порту 8080

3. **Проверьте статус:**
   ```bash
   docker-compose ps
   ```

4. **Проверьте логи:**
   ```bash
   docker-compose logs -f api
   ```

5. **Протестируйте API:**
   ```bash
   curl http://localhost:8080/team/get?team_name=test
   ```


### Архитектура

Проект следует принципам чистой архитектуры с разделением на слои:

- **Handler** - обработка HTTP запросов, валидация, преобразование DTO
- **Service** - бизнес-логика, правила домена
- **Repository** - работа с базой данных, SQL запросы

Каждый слой имеет свои DTO и не знает о структурах вышестоящих слоёв.

### Makefile команды

```bash
make help              # Показать справку по командам
make deps              # Установить зависимости
make build             # Собрать бинарник
make run               # Запустить сервис локально
make test              # Запустить тесты
make test-coverage     # Запустить тесты с покрытием
make fmt               # Форматировать код
make clean             # Очистить сгенерированные файлы

# Docker команды
make docker-build      # Собрать Docker образы
make docker-up         # Запустить сервисы через Docker Compose
make docker-down       # Остановить сервисы
make docker-logs       # Показать логи
make docker-restart    # Перезапустить сервисы
```

### Тестирование

#### Запуск E2E тестов

E2E тесты используют реальную БД и проверяют полные сценарии работы API:

```bash
# Убедитесь, что PostgreSQL запущен
make test

# Или напрямую
go test ./tests/e2e -v
```

Тесты автоматически:
- Подключаются к БД (по умолчанию `PReviewer`)
- Применяют миграции
- Очищают данные перед каждым тестом
- Запускают HTTP сервер для тестирования

#### Переменные окружения для тестов

```bash
export TEST_DATABASE_URL="postgres://user:password@localhost:5432/test_db?sslmode=disable"
go test ./tests/e2e -v
```

### Конфигурация

Основной конфигурационный файл: `configs/apiserver.toml`

```toml
bind_addr = ":8080"
log_level = "debug"

[store]
database_url = "postgres://appuser:secret@localhost:5432/PReviewer?sslmode=disable"
```

**Переменные окружения:**

- `DATABASE_URL` - переопределяет `database_url` из конфига (приоритет над файлом)

### База данных

#### Миграции

Миграции находятся в `db/migrations/` и применяются автоматически при запуске PostgreSQL через Docker.

**Структура БД:**

- `teams` - команды
- `users` - пользователи (связь с командами)
- `pullrequests` - PR'ы (связь с авторами и ревьюверами)

#### Подключение к БД

**Через Docker:**
```bash
docker-compose exec db psql -U appuser -d PReviewer
```

**Локально:**
```bash
psql -U appuser -d PReviewer -h localhost
```

##  Docker

### Сборка образов

```bash
docker-compose build
```

### Запуск

```bash
# В фоновом режиме
docker-compose up -d

# С логами
docker-compose up
```

### Остановка

```bash
# Остановить сервисы
docker-compose down

# Остановить и удалить volumes (очистить БД)
docker-compose down -v
```

### Логи

```bash
# Все сервисы
docker-compose logs -f

# Только API
docker-compose logs -f api

# Только БД
docker-compose logs -f db
```

### Пересборка после изменений

```bash
docker-compose build --no-cache api
docker-compose up -d
```
