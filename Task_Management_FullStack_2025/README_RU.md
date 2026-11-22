# Task Management Full-Stack 2025 — Restored Edition

**Проект №3** из истории учебных разработок: Go backend + React frontend, создание/изменение/удаление задач, маршруты, сроки, статусы.

**Важно:** это новая совместимая реконструкция по старому заданию и открытой информации о проектах команды, а НЕ копия файлов из GitHub коллеги. Первоначальные публичные репозитории: 
- Backend: https://github.com/Bekzhanizb/TaskManagerBackend
- Frontend: https://github.com/Bekzhanizb/TaskManagerFrontend

## Быстрый запуск на твоём Windows 10 (без установок)

1. Распакуй ZIP **целиком**. Не запускай программу прямо внутри ZIP.
2. Дважды нажми `START_WINDOWS.bat`.
3. Когда появится `Task Manager Full-Stack started`, открой Chrome: **http://localhost:8082**.
4. Для остановки сервера нажми `Ctrl+C` в CMD.
5. Созданные задачи сохраняются в `backend/data/tasks.json`, даже после перезапуска.

Можно запускать и из CMD (команды по отдельности):

```cmd
cd /d C:\Users\Acer\Downloads\Task_Management_FullStack_2025
START_WINDOWS.bat
```

Не нужна Java, Go, Node.js или Docker! В ZIP входит `TaskManagement.exe` — собранный Go сервер под Windows x64, со встроенным React-интерфейсом. Быстрый режим слушает только **localhost** и предназначен для одного пользователя, не для публичного размещения в интернете. Режим быстрого запуска хранит задачи **в JSON-файле**, **не в GORM/PostgreSQL**. Это удобно для мгновенной проверки интерфейса.

## Полный учебный стек Go + GORM + PostgreSQL

1. Установи и запусти **Docker Desktop** (Windows, WSL 2 backend).
2. Сначала останови `START_WINDOWS.bat`, потому что обе версии используют порт `8082`.
3. Запусти `START_DOCKER_WINDOWS.bat` (или выполни `docker compose up --build`).
4. Открой **http://localhost:8082**.
5. При первом запуске Docker скачивает PostgreSQL, Go image и модули Go (GORM).
6. Данные сохраняются в Docker volume `pg_data`. Перезапуск `docker compose up` их не удаляет. `docker compose down -v` **удалит** эту базу.

Этот режим использует настоящий **GORM** (`backend/storage_postgres.go`) и **PostgreSQL 16**. SQL база и JSON-файл портативного режима независимы друг от друга: задачи не переносятся автоматически.

## Функции

- Список задач / Dashboard и прогресс.
- Создание, редактирование, удаление.
- Название, описание, срок, приоритет, статус, done.
- Поиск, фильтры статуса и приоритета.
- Kanban-доска с перетаскиванием мышкой.
- Статистика, просроченные задачи, CSV экспорт.
- Демонстрационные задачи (добавляются только по кнопке).
- REST API и сохранение данных между запусками.

## Технологии

| Часть | Реализация |
|---|---|
| Backend | Go (net/http), JSON REST API |
| Frontend | React 16 UMD, JavaScript, адаптивный CSS |
| Быстрый запуск | Стандартный Go / локальный JSON |
| Полный Docker | Go + GORM + PostgreSQL 16 |
| Интерфейс/API | `http://localhost:8082` |
| Windows | Встроенный `.exe` для AMD64 |

Исходная командная работа в 2025 году использовала Gorilla Mux, Axios, React Router и Bootstrap. Для стабильного полностью автономного запуска в реконструкции используется стандартный `net/http`, локальный React и CSS без CDN. REST пути и возможности задач сохранены.

## API

| Метод | Адрес | Назначение |
|---|---|---|
| GET | `/health` | Работа сервера |
| GET | `/api/tasks` | Все задачи |
| POST | `/api/tasks` | Создать задачу |
| GET | `/api/tasks/{id}` | Детали задачи |
| PUT | `/api/tasks/{id}` | Изменить задачу |
| DELETE | `/api/tasks/{id}` | Удалить задачу |
| GET | `/api/stats` | Статистика |
| GET | `/api/export.csv` | Скачать CSV |

Пример JSON:

```json
{"title":"Finish homework","description":"Task management demo","due_date":"2026-12-01","priority":"high","status":"todo"}
```

Поддерживается `dueDate` (старое название в React) как псевдоним `due_date`.

## Проверка

- Открой сайт и создай задачу.
- Отметь её выполненной; перейди в Kanban.
- Останови сервер, снова запусти `START_WINDOWS.bat`, убедись, что задача сохранилась.
- Во второй CMD запусти `TEST_API_WINDOWS.bat`: проверяет HTTP health и CRUD.
- Для Go-тестов при установленном Go выполни: `cd backend` → `go test ./...`.

**Выполненная проверка при подготовке ZIP:** стандартный Go backend успешно компилировался в Windows x64, прошёл тесты API/JSON-персистентности/валидации/CORS; созданный ZIP проверен на целостность. Режим GORM/PostgreSQL Docker в этой среде не запускался, потому что недоступны загрузка модулей Go и контейнерный движок. Итоговый запуск проверяй в Docker Desktop.

## GitHub: Zhenis455

Создай пустой репозиторий `task-management-fullstack-2025` на https://github.com/new. В той же папке проекта открой CMD и вводи **по одной строке**:

```cmd
git init
git add .
git commit -m "Task Management Full-Stack 2025 - restored"
git branch -M main
git remote add origin https://github.com/Zhenis455/task-management-fullstack-2025.git
git push -u origin main
```

Если `origin` уже существует: `git remote set-url origin https://github.com/Zhenis455/task-management-fullstack-2025.git`.

`backend/data/tasks.json` игнорируется Git, чтобы не публиковать персональные данные. **Не переписывай историю оригинальных репозиториев товарища, загружай эту новую реализацию в свой отдельный репозиторий.**

## Структура

```
START_WINDOWS.bat             # один запуск без установки
TaskManagement.exe           # Go Windows executable (frontend embedded)
TEST_API_WINDOWS.bat          # проверка API
OPEN_CODE_WINDOWS.bat        # открыть код в VS Code
START_DOCKER_WINDOWS.bat      # полный Go/GORM/PostgreSQL режим
compose.yaml
backend/
  server.go                  # маршруты REST и React server
  model.go                   # поля, валидаторы
  storage_file.go            # локальный режим
  storage_postgres.go        # GORM режим, build tag postgres
  server_test.go             # автоматические тесты
  web/                       # встроенный React frontend
frontend/src/app.js          # исходник React
scripts/verify-api.ps1
```
