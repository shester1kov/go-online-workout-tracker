# Отчет 4: Настройка Nginx как Reverse Proxy для Go-приложения с элементами безопасности

## Настройка Nginx как Reverse Proxy для Go-приложения с элементами безопасности

### Цель

Настроить Nginx как Reverse Proxy для Go-приложения, реализовать балансировку нагрузки и применить базовые меры защиты.

## 1. Установка и настройка Nginx в Docker

- Использован официальный образ `nginx:1.27.4-alpine`.
- Конфигурация Nginx монтируется из локальной директории `./nginx/`.
- Файл `nginx.conf`:

```nginx
user nginx;
worker_processes auto;

error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;
    server_tokens off;

    limit_req_zone $binary_remote_addr zone=one:10m rate=10r/s;
    access_log /var/log/nginx/access.log;
    sendfile on;
    keepalive_timeout 65;

    include /etc/nginx/conf.d/*.conf;
}
```

## 2. Балансировка нагрузки

- Запущено три экземпляра Go-приложения (`backend`, `backend2`, `backend3`).
- Использован алгоритм балансировки `least_conn` (минимальное число соединений).
- Файл `conf.d/app.conf`:

```nginx
upstream backend {
    least_conn;
    server backend:8080;
    server backend2:8081;
    server backend3:8082;
}

server {
    listen 80;
    server_name localhost;

    client_max_body_size 2M;

    location / {
        limit_req zone=one burst=20 nodelay;
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /api/v1/login {
        limit_req zone=one burst=5 nodelay;
        proxy_pass http://backend;
    }

    if ($request_method !~ ^(GET|POST|PUT|DELETE|HEAD)$) {
        return 403;
    }
}
```

## 3. Настройка HTTPS (не реализовано)

- Настройка HTTPS не выполнена, так как отсутствует домен.
- Планируется использование Let's Encrypt через Certbot.

## 4. Меры безопасности

- Ограничение числа запросов:

```nginx
limit_req_zone $binary_remote_addr zone=one:10m rate=10r/s;
```

- Ограничение методов запросов:

```nginx
if ($request_method !~ ^(GET|POST|PUT|DELETE|HEAD)$) {
    return 403;
}
```

- Ограничение доступа к `/api/v1/login` (защита от брутфорса):

```nginx
location /api/v1/login {
    limit_req zone=one burst=5 nodelay;
    proxy_pass http://backend;
}
```

- Отключение информации о сервере:

```nginx
server_tokens off;
```

---

## 5. Логирование и мониторинг

- Логирование запросов:

```nginx
access_log /var/log/nginx/access.log;
error_log /var/log/nginx/error.log warn;
```

- Установлен и настроен Fail2Ban в Docker:

```yaml
fail2ban:
  image: crazymax/fail2ban:latest
  restart: always
  network_mode: "host"
  cap_add:
    - NET_ADMIN
    - NET_RAW
  volumes:
    - ./fail2ban:/data
    - ./nginx/logs:/var/log/nginx:ro
  environment:
    - TZ=Europe/Moscow
```

- Конфигурация `/etc/fail2ban/jail.local`:

```ini
[nginx-http-auth]
enabled = true
filter = nginx-http-auth
action = iptables-multiport[name=HTTP, port="80,443", protocol=tcp]
logpath = /var/log/nginx/error.log
bantime = 600
maxretry = 5
```

---

## Результаты работы

- Nginx работает как Reverse Proxy в контейнере.
- Реализована балансировка нагрузки между тремя экземплярами Go-приложения.
- Настроены базовые меры защиты: ограничение запросов, скрытие информации о сервере.
- Установлен и настроен Fail2Ban для защиты от атак.
- Настройка HTTPS не выполнена (отсутствует домен).
