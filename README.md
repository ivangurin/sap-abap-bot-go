# SAP ABAP Bot

Телеграм-бот, который отвечает на вопросы по SAP и ABAP с помощью ИИ. Создан для группы
[SAP ABAP](https://t.me/sapabap).

## Как работает

Бот подключается к ИИ-провайдеру и отвечает на обращения:

- **В группе** — отвечает, если его упомянули через `@<bot>` или если он уже участвует в
  треде (продолжает диалог). Контекст треда хранится 24 часа.
- **В личке** — отвечает только администраторам (`SAP_ABAP_BOT_ADMIN_USER_IDS`).

Если к боту обратились, ответ он даёт всегда. Ответ форматируется под Telegram (жирный
текст и блоки кода).

## Провайдеры ИИ

Поддерживаются два провайдера, оба на официальных SDK. Выбор — переменной
`SAP_ABAP_BOT_AI_PROVIDER`:

- `github` (по умолчанию) — [GitHub Models](https://github.com/marketplace/models),
  OpenAI-совместимый API (`github.com/openai/openai-go`).
- `anthropic` — Anthropic API (`github.com/anthropics/anthropic-sdk-go`).

У каждого провайдера настраиваются базовый url, токен и модель (см. таблицу ниже).

## Конфигурация

Настройка через переменные окружения (можно через файл `.env`, пример — `.env-example`):

| Переменная | Описание | По умолчанию |
| --- | --- | --- |
| `SAP_ABAP_BOT_TOKEN` | Токен телеграм-бота | — |
| `SAP_ABAP_BOT_ADMIN_USER_IDS` | ID администраторов (через запятую), кому бот отвечает в личке | — |
| `SAP_ABAP_BOT_ALLOWED_CHAT_IDS` | ID разрешённых чатов (через запятую); пусто — без ограничений | — |
| `SAP_ABAP_BOT_AI_PROVIDER` | Провайдер ИИ: `github` или `anthropic` | `github` |
| `SAP_ABAP_BOT_GITHUB_BASE_URL` | Базовый url GitHub Models | `https://models.inference.ai.azure.com` |
| `SAP_ABAP_BOT_GITHUB_TOKEN` | Токен GitHub-аккаунта | — |
| `SAP_ABAP_BOT_GITHUB_AI_MODEL` | Модель GitHub Models | `gpt-4.1` |
| `SAP_ABAP_BOT_ANTHROPIC_BASE_URL` | Базовый url Anthropic API (опционально, для своего шлюза) | — |
| `SAP_ABAP_BOT_ANTHROPIC_TOKEN` | Токен Anthropic API | — |
| `SAP_ABAP_BOT_ANTHROPIC_AI_MODEL` | Модель Anthropic | `claude-opus-4-8` |
| `SAP_ABAP_BOT_LOG_LEVEL` | Уровень логирования (`debug`/`info`/`warn`/`error`) | `debug` |
| `SAP_ABAP_BOT_LOG_FILE` | Файл для логов; пусто — вывод в stdout | — |

Заполнять нужно токены только того провайдера, который выбран в `SAP_ABAP_BOT_AI_PROVIDER`.

## Запуск

Локально:

```bash
make run
```

Через Docker Compose (образ `ghcr.io/ivangurin/sap-abap-bot-go`):

```bash
make pull-and-run
```

## Разработка

- `make lint` — линтер (golangci-lint).
- `make genmock` — перегенерация моков (mockery).
- `make test` — тесты (`go test -race ./...`).

Технологии: Go 1.26, [go-telegram/bot](https://github.com/go-telegram/bot),
официальные SDK OpenAI и Anthropic.

Устройство проекта и правила его поддержки описаны в [CLAUDE.md](./CLAUDE.md).
