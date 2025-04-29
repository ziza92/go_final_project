# Планировщик задач в Веб-Сервере

## Описание
Этот проект представляет собой TODO LIST, с поддержками расписаний и повторений задач.

## Инструкция по установке локально
  - `TODO_PORT`: (Опционально) Порт для сервера по умолчанию: `7540`.
  - `TODO_DBFILE`: (Опционально) Располажение базы данных по умолчанию в директории проекта.
  - `TODO_PASSWORD`: Если переменная среды установлена, требуется аутентификация. Иначе, аутентификация отключена.
  ```bash
   git clone https://github.com/ziza92/go_final_project.git
   cd go_final_project
   go mod tidy
   go run .
   ```

## Инструкция по установке с использованием Docker
  ```bash
   git clone https://github.com/ziza92/go_final_project.git
   cd go_final_project
   docker compose up -d
   ```

### Запуск сервера
Откройте браузер и перейдите по адресам:
   - `http://localhost:7540/` — главная страница планировщика.
   - `http://localhost:7540/login.html` — страница входа (если установлен `TODO_PASSWORD`).

   **Пароль TODO_PASSWORD по умолчанию установлен в вариации с Docker использованием.**<br/>
   Его значение - ```primite-plz```