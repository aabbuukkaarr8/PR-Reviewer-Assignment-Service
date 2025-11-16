FROM postgres:16.6-alpine3.20

COPY db/migrations/*.up.sql /docker-entrypoint-initdb.d/
COPY db/insert/*.sql /docker-entrypoint-initdb.d/

WORKDIR /docker-entrypoint-initdb.d
