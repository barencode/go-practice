# Go practice

Изучаю Go: упражнения, повторение и практика Git

Первый коммит создан. Учусь сохранять изменения осознанно.

- Запуск базы: docker compose up -d.
- Первоначальное создание таблиц: две команды применения SQL-файлов, сначала 001, затем 002.
    - docker compose exec -T db psql -U mining -d mining < migrations/001_create_mines.up.sql
    - docker compose exec -T db psql -U mining -d mining < migrations/002_create_miners.up.sql
- Начальное наполнение: docker compose exec -T db psql -U mining -d mining -v ON_ERROR_STOP=1 < scripts/seed.sql
- Данные загружаются при запуске сохранение данных по HTTP TBD
- Подключение: docker compose exec db psql -U mining -d mining.
- Остановка: docker compose down, при которой данные в именованном томе сохраняются.

- Запуск сервера: DATABASE_URL='postgres://mining:mining_dev@127.0.0.1:5432/mining?sslmode=disable' go run .