FROM golang:latest as build

WORKDIR /workdir

COPY . /workdir

RUN apt update && apt install -y make zip curl
RUN make build
RUN make get-front


FROM ubuntu:latest

WORKDIR /oj-lab-platform

COPY --from=build /workdir/bin/web_server /usr/local/bin/web_server
COPY --from=build /workdir/frontend/dist ./frontend_dist

COPY config.toml ./config.toml

ENV OJ_LAB_SERVICE_ENV='production'
ENV DATABASE_DSN='user=postgres password=postgres host=host.docker.internal port=5432 dbname=oj_lab sslmode=disable TimeZone=Asia/Shanghai'
ENV REDIS_HOSTS='["host.docker.internal:6379"]'
ENV MINIO_ENDPOINT='http://host.docker.internal:9000'

EXPOSE 8080
CMD [ "web_server" ]