## Файлы и их назначения
### cmd/server/main.go: 
Основной файл для запуска HTTP-сервера.
### config/config.go: 
Код для загрузки конфигурационных данных из файла.
### config/config.yaml: 
YAML файл с конфигурацией приложения, включая параметры базы данных и сервера.
### internal/api/handler.go: 
Обработчики HTTP-запросов для создания, получения и остановки команд.
### internal/db/db.go: 
Инициализация базы данных и функции для выполнения SQL-запросов.
### internal/models/command.go: 
Модель данных для представления команд.
### test/api_test.go: 
Тесты для API, проверяющие функциональность создания и получения команд.
### Dockerfile: 
Скрипт для сборки Docker-образа приложения.
### docker-compose.yml: 
Конфигурация Docker Compose для запуска приложения и базы данных PostgreSQL.

## Запуск проекта

### Создайте базу данных и пользователя в PostgreSQL:

````
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo -i -u postgres
psql
CREATE DATABASE command_db;
CREATE USER postgres WITH ENCRYPTED PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON DATABASE command_db TO postgres;
````

### Склонируйте репозиторий и перейдите в директорию проекта:

````
git clone https://github.com/beValed/bash_command_api.git
cd bash_command_api
````

### Запустите проект с Docker:

````
docker-compose up --build
````

### API будет доступно по адресу: http://localhost:8080

