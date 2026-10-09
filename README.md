# github-activity

Небольшое CLI-приложение на Go, которое получает последнюю публичную активность
пользователя GitHub через REST API и выводит её в консоль в человекочитаемом виде.

Идея проекта взята с [roadmap.sh](https://roadmap.sh/projects/github-user-activity).

## Требования

- Go 1.24+

## Установка и запуск

```bash
go build -o github-activity ./cmd/github-activity
./github-activity <username>
```

Или сразу через `go run`:

```bash
go run ./cmd/github-activity torvalds
```

## Пример вывода

```
Pushed to master in torvalds/linux
Created branch preserve-identity in gaearon/react
Opened pull request #35 in sindresorhus/eslint-cssicorn
Merged pull request #48 in torvalds/GuitarPedal
Published release v2.2.0 in sindresorhus/ky
Starred some/repository
```

В терминале основной глагол действия подсвечивается зелёным, а название
репозитория — жёлтым. Цвет автоматически отключается при перенаправлении вывода
в файл или канал, а также при установленной переменной окружения `NO_COLOR`.

## Как это работает

Программа обращается к эндпоинту
[`GET /users/<username>/events`](https://docs.github.com/en/rest/activity/events)
и превращает сырые JSON-события в понятные фразы.

> **Примечание:** публичный фид событий не содержит количество коммитов в
> `PushEvent` (в payload есть только `ref`, `before`, `head`), поэтому для пуша
> выводится название ветки: `Pushed to master in torvalds/linux`. Если число
> коммитов всё же присутствует в ответе, оно используется: `Pushed 3 commits to ...`.

## Структура проекта

```
cmd/github-activity/   — точка входа, разбор аргументов, вывод
internal/model/        — доменные типы: Event, Repo, EventType
internal/client/       — клиент GitHub API, получение событий
internal/format/       — форматирование событий в текст (с цветом)
```

## Разработка

```bash
gofmt -l .
go vet ./...
go test ./...
```
