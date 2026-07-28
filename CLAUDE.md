# CLAUDE.md

Памятка для Claude Code по проекту `sap-abap-bot-go`.

## ⚠️ Поддержка этого файла

**При ЛЮБЫХ изменениях в проекте обновляй этот файл.** Добавил провайдера, поменял
архитектуру, добавил env-переменную, изменил поведение бота, команды или конвенции —
сразу отрази это в CLAUDE.md в том же изменении. Файл должен всегда соответствовать
текущему состоянию кода.

## Обзор

Телеграм-бот, который отвечает на вопросы по SAP и ABAP с помощью ИИ. Бот вызывается,
когда к нему обращаются, передаёт диалог выбранному ИИ-провайдеру и публикует ответ.

- Язык: Go 1.26, модуль `bot`.
- Точка входа: `cmd/bot/main.go` → `internal/app/bot/app.go`.

## Архитектура

Ручной DI-контейнер (`internal/service/provider`) лениво создаёт клиентов и сервисы.

- `internal/config` — конфиг из env (`caarlos0/env` + `.env` через `godotenv`).
  Здесь же хранится системный промпт (`App.SystemPrompt`).
- `internal/model` — общие доменные типы: `Thread`/`ThreadMessage` и
  провайдеро-независимый `ChatMessage` (роль + текст).
- `internal/clients/github` — клиент GitHub Models через **официальный OpenAI SDK**
  (`github.com/openai/openai-go`), эндпоинт OpenAI-совместимый.
- `internal/clients/anthropic` — клиент Anthropic через **официальный SDK**
  (`github.com/anthropics/anthropic-sdk-go`).
- `internal/service/agent` — провайдеро-независимый агент: собирает диалог
  (`[]*model.ChatMessage`) и вызывает `Client.Ask`, оборачивая результат в `[]*Answer`.
- `internal/service/bot` — телеграм-слой (`go-telegram/bot`): хендлер, гейтинг,
  хранение тредов в памяти.
- `internal/service/provider` — сборка зависимостей; `getAIClient()` выбирает клиента
  по `config.App.AIProvider`.
- `internal/pkg` — logger, closer, utils.

## ИИ-провайдеры

Оба клиента реализуют один интерфейс (структурно):

```go
Ask(ctx context.Context, systemPrompt string, messages []*model.ChatMessage) ([]string, error)
```

- Выбор провайдера — env `SAP_ABAP_BOT_AI_PROVIDER`: `github` (по умолчанию) или `anthropic`
  (константы `config.AIProviderGitHub` / `config.AIProviderAnthropic`).
- У обоих настраиваются `base_url`, `token`, `model` (см. env ниже).
- Инструментов (tool use / function calling) **нет намеренно**: бот всегда должен
  ответить, поэтому ответ берётся из обычного текстового ответа модели.

Чтобы добавить нового провайдера: создать пакет в `internal/clients/<name>` с методом
`Ask`, добавить геттер в `provider/clients.go`, ветку в `provider/services.go`
(`getAIClient`) и константу/поле конфига.

## Конфигурация (env)

- `SAP_ABAP_BOT_TOKEN` — токен телеграм-бота.
- `SAP_ABAP_BOT_ADMIN_USER_IDS` — id админов (для лички).
- `SAP_ABAP_BOT_ALLOWED_CHAT_IDS` — id разрешённых чатов.
- `SAP_ABAP_BOT_AI_PROVIDER` — `github` | `anthropic`.
- `SAP_ABAP_BOT_GITHUB_BASE_URL` / `_TOKEN` / `_AI_MODEL`.
- `SAP_ABAP_BOT_ANTHROPIC_BASE_URL` / `_TOKEN` / `_AI_MODEL`.
- `SAP_ABAP_BOT_LOG_LEVEL`, `SAP_ABAP_BOT_LOG_FILE`.

При изменении набора env-переменных обнови также `.env-example`, `docker-compose.yml`
и `README.md`.

## Команды (Makefile)

- `make run` — запуск (`go run ./cmd/bot`).
- `make lint` — golangci-lint с `--fix` и `--new-from-rev=master`.
- `make genmock` — перегенерация моков (`go tool mockery`).
- `make test` — `go test -race ./...` (тестов пока нет).

## Поведение бота

- Гейтинг в `service/bot/dafault_handler.go`: в личке отвечает только админам; в группе —
  только если тред уже отслеживается (бот в переписке) или в сообщении есть `@<bot>`.
- Если бота вызвали — ответ обязателен (это отражено в системном промпте).
- Треды хранятся в памяти (`service/bot/thread.go`), чистятся раз в час, живут 24 часа.
- Форматирование ответа — Telegram MarkdownV2. Хендлер экранирует спецсимволы и надёжно
  пропускает только жирный текст `**...**` и блоки кода в тройных кавычках; системный
  промпт ограничивает модель именно этим форматом.

## Конвенции

- Комментарии и тексты — на русском, как в остальном коде.
- Мок GitHub-клиента генерируется mockery (см. `.mockery.yaml`); после изменения
  интерфейса `github.Client` запускай `make genmock`.
- Перед завершением задачи прогоняй `go build ./...`, `go vet ./...` и `make lint`.
