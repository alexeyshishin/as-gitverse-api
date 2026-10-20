# as-gitverse-api

CLI для массового управления репозиториями [Gitverse](https://gitverse.ru) через
публичный API. Авторизация — Personal Access Token.

## Установка

```
go build -o gitverse ./cmd/gitverse
```

## Токен

Приоритет: флаг `--token` > переменная `GITVERSE_TOKEN` > файл
`~/Library/Application Support/gitverse/credentials` (macOS) /
`~/.config/gitverse/credentials` (Linux), права 0600.

```
export GITVERSE_TOKEN=xxxxxxxxxxxxxxxx
```

## Глобальные флаги

| Флаг | Описание |
|---|---|
| `--token` | PAT (переопределяет `GITVERSE_TOKEN` и файл креденшлов) |
| `--json` | вывод в JSON вместо таблицы |
| `--concurrency` | макс. параллельных запросов в bulk-операциях (по умолчанию 5) |
| `--verbose` | подробное логирование (токен никогда не попадает в лог) |

## Команды

### `repo` — репозитории

| Команда | Описание |
|---|---|
| `repo create --names/--file [--private]` | создать один или несколько репозиториев (bulk) |
| `repo list [--org <org>]` | список репозиториев текущего пользователя или организации |
| `repo collaborators list --owner --repo` | список коллабораторов репозитория |
| `repo collaborators add --owner --repo --names/--file --permission` | добавить коллабораторов (bulk; `--permission`: `read`/`write`/`admin`) |
| `repo commit-stats --owner --repo [--since] [--until]` | статистика коммитов по авторам |

```
# создать несколько публичных репо за один вызов
gitverse repo create --names foo,bar,baz

# создать приватные репо из файла (по одному имени на строку)
gitverse repo create --file repos.txt --private

# список своих репозиториев
gitverse repo list

# список репозиториев организации
gitverse repo list --org studentlabs

# список коллабораторов репозитория
gitverse repo collaborators list --owner studentlabs --repo slide_backend

# добавить коллабораторов с правом write
gitverse repo collaborators add --owner studentlabs --repo slide_backend --names alice,bob --permission write

# статистика коммитов по авторам
gitverse repo commit-stats --owner studentlabs --repo slide_backend

# статистика коммитов за период (--since/--until — даты YYYY-MM-DD, --until включительно)
gitverse repo commit-stats --owner studentlabs --repo slide_backend --since 2026-08-01 --until 2026-08-28
```

## Вывод в JSON

Любую команду можно вывести в JSON вместо таблицы:

```
gitverse repo list --org studentlabs --json
```

## Что осознанно не реализовано

- Любые `delete`-команды для репозиториев — запрещено жёстким правилом проекта
  (удаление — только вручную через UI Gitverse).
- `collaborators remove` и `org members` (полный листинг участников
  организации) — проверено на публичном API: таких эндпоинтов не существует
  (`GET /orgs/{org}/members` без username → `400`, есть только
  `GET /orgs/{org}/members/{username}` — проверка членства одного юзера).
- Вебхуки (`webhook`) и команды организации (`team`) — исключены из скоупа
  проекта, не реализуются.
