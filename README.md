# Journal

Электронный журнал для семейных классов. Проект состоит из React-клиента, Go API и PostgreSQL.

## Требования

Для локального запуска нужны:

- Docker с Docker Compose;
- Go 1.25 или новее;
- Node.js 20.19+ или 22.12+;
- pnpm;
- Make.

## Быстрый запуск

Все команды ниже выполняются из корня проекта, если не указано иное.

### 1. Запустить PostgreSQL

```shell
docker compose up -d
docker compose ps
```

PostgreSQL будет доступен на `127.0.0.1:5433`. Данные сохраняются в Docker volume `journal_postgres_data`.

### 2. Подготовить и запустить сервер

Сначала установите зафиксированные версии Goose и генератора OpenAPI:

```shell
cd server
make deps
make migrate-up
```

Запустите API. Для локального HTTP обязательно отключите флаг Secure у cookie.

PowerShell:

```powershell
$env:COOKIE_SECURE = "false"
make run
```

Bash:

```bash
export COOKIE_SECURE=false
make run
```

По умолчанию сервер запускается на [http://127.0.0.1:8080](http://127.0.0.1:8080), а Swagger UI доступен на [http://127.0.0.1:8080/docs](http://127.0.0.1:8080/docs).

Сервер использует следующую строку подключения по умолчанию:

```text
postgres://postgres:password@127.0.0.1:5433/journal?sslmode=disable
```

Чтобы использовать другую базу данных или адрес сервера, задайте переменные `DATABASE_URL` и `HTTP_ADDR` перед запуском. `SHUTDOWN_TIMEOUT` задаёт срок корректного завершения (по умолчанию `10s`). Некорректные значения конфигурации останавливают запуск с понятной ошибкой.

### 3. Запустить клиент

Откройте второй терминал:

```shell
cd client
pnpm install
pnpm dev
```

Клиент будет доступен на [http://127.0.0.1:5173](http://127.0.0.1:5173). Vite автоматически проксирует запросы `/api` на `http://127.0.0.1:8080`.

## Остановка

Остановите клиент и сервер сочетанием `Ctrl+C`, затем остановите PostgreSQL:

```shell
docker compose down
```

Команда сохраняет данные базы. Для обычной остановки не используйте флаг `-v`, поскольку он удаляет volume с данными PostgreSQL.

## Проверки

Сервер:

```shell
cd server
make test
make lint
make build
make generate
make migrate-validate
make integration-test
```

Клиент:

```shell
cd client
pnpm test
pnpm run lint
pnpm run build
```

## Полезные команды

Проверить состояние миграций:

```shell
cd server
make migrate-status
```

Применить новые миграции:

```shell
cd server
make migrate-up
```

Пересоздать Go-типы после изменения OpenAPI:

```shell
cd server
make generate
```

## Возможные проблемы

- Если сервер не подключается к PostgreSQL, убедитесь, что Docker запущен и контейнер имеет статус `healthy` в `docker compose ps`.
- Если браузер после входа снова показывает форму авторизации, проверьте, что сервер запущен с `COOKIE_SECURE=false`.
- Если Vite, Vitest или ESLint сообщают об ошибках синтаксиса Node.js, проверьте версию командой `node --version`.
- Команды `make migrate-up` и `make run` нужно выполнять из каталога `server`.
- Каждый HTTP-ответ содержит `X-Request-ID`; тот же идентификатор присутствует в структурированном JSON-логе сервера. Тела запросов, query-параметры, cookie и заголовки авторизации не логируются.
