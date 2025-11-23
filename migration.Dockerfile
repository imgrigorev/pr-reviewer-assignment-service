#FROM alpine:3.13
#
#RUN apk update && \
#    apk upgrade && \
#    apk add --no-cache bash curl postgresql-client && \
#    rm -rf /var/cache/apk/*
#
## Скачиваем бинарник migrate
#RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz -o /tmp/migrate.tar.gz && \
#    tar -xzf /tmp/migrate.tar.gz -C /bin/ && \
#    chmod +x /bin/migrate && \
#    rm /tmp/migrate.tar.gz
#
#WORKDIR /root
#
## Копируем миграции и скрипты
#COPY migrations/*.sql ./migrations/
#COPY migration.sh .
#COPY .env .
#
#RUN chmod +x migration.sh
#
#ENTRYPOINT ["bash", "migration.sh"]
FROM alpine:3.13

RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash curl postgresql-client && \
    rm -rf /var/cache/apk/*

# Скачиваем правильную версию migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz -C /bin/ && \
    chmod +x /bin/migrate

WORKDIR /root

# Копируем миграции и скрипты
COPY migrations/*.sql ./migrations/
COPY migration.sh .
COPY .env .

RUN chmod +x migration.sh

ENTRYPOINT ["bash", "migration.sh"]