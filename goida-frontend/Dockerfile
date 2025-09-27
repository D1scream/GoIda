# Используем nginx для раздачи статических файлов
FROM nginx:alpine

# Копируем файлы фронтенда
COPY index.html /usr/share/nginx/html/
COPY app.js /usr/share/nginx/html/
COPY styles.css /usr/share/nginx/html/

# Настраиваем nginx для SPA
RUN echo 'server { \
    listen 80; \
    server_name localhost; \
    root /usr/share/nginx/html; \
    index index.html; \
    \
    # CORS заголовки для работы с API \
    add_header Access-Control-Allow-Origin *; \
    add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS"; \
    add_header Access-Control-Allow-Headers "Content-Type, Authorization"; \
    \
    # Обработка preflight запросов \
    location / { \
        if ($request_method = OPTIONS) { \
            add_header Access-Control-Allow-Origin *; \
            add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS"; \
            add_header Access-Control-Allow-Headers "Content-Type, Authorization"; \
            add_header Content-Length 0; \
            add_header Content-Type text/plain; \
            return 200; \
        } \
        try_files $uri $uri/ /index.html; \
    } \
    \
    # Кэширование статических файлов \
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ { \
        expires 1y; \
        add_header Cache-Control "public, immutable"; \
    } \
}' > /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
