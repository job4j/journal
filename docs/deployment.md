# Развёртывание Journal

Journal публикуется на основном домене по адресу
[https://job4j.ru/journal/](https://job4j.ru/journal/).

Jenkins собирает сервер и клиент, применяет миграции, обновляет systemd-сервис
`journal-api` и распаковывает клиент в `/var/www/job4j.ru/journal`.

## Подготовка сервера

1. Создайте базу данных и отдельного пользователя PostgreSQL.
2. Создайте файл `/etc/journal/server.env`:

   ```text
   DATABASE_URL=postgres://journal:CHANGE_ME@127.0.0.1:5432/journal?sslmode=disable
   HTTP_ADDR=127.0.0.1:9082
   COOKIE_SECURE=true
   SHUTDOWN_TIMEOUT=10s
   ```

3. Ограничьте доступ к файлу:

   ```shell
   sudo chown root:deploy /etc/journal/server.env
   sudo chmod 640 /etc/journal/server.env
   ```

4. Создайте каталог клиента и ссылку внутри корневого сайта:

   ```shell
   sudo mkdir -p /opt/journal/client
   sudo ln -s /opt/journal/client /var/www/job4j.ru/journal
   ```

   Если ссылка уже существует, не создавайте её повторно. Jenkins пишет в
   каталог через `sudo`.
5. Добавьте пользователю `deploy` права `sudo` без пароля для команд,
   используемых Jenkins: `mkdir`, `chown`, `install`, `tar`, `tee`
   и `systemctl`.
6. В Jenkins создайте Pipeline из репозитория и укажите `Jenkinsfile`.
   Credential `job4j-pixel-deploy-ssh` должен содержать SSH-ключ пользователя
   `deploy`. Если используется другое имя credential, измените
   `SSH_CREDENTIALS_ID` в `Jenkinsfile`.
7. Убедитесь, что из контейнера Jenkins адрес `host.docker.internal`
   открывает SSH хоста. При другой схеме измените `DEPLOY_HOST`.
8. Проверьте nginx и примените конфигурацию:

   ```shell
   sudo nginx -t
   sudo systemctl reload nginx
   ```

После первого запуска Jenkins сам установит systemd unit, применит миграции и
запустит API. Создавать таблицы вручную не нужно: вручную создаётся только база
данных и пользователь PostgreSQL.

## Проверка

```shell
curl -I https://job4j.ru/journal/
curl -i https://job4j.ru/journal/api/v1/me
sudo systemctl status journal-api
sudo journalctl -u journal-api -n 100 --no-pager
```

Неавторизованный запрос `/journal/api/v1/me` должен вернуть `401`, а не
`404` или `502`.
