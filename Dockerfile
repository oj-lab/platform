FROM golang:latest as build

WORKDIR /workdir

COPY . /workdir

RUN apt update && apt install -y make zip curl
RUN make build
RUN make get-front


FROM ubuntu:latest

RUN apt update && apt install -y ca-certificates curl
RUN update-ca-certificates

WORKDIR /platform

COPY --from=build /workdir/bin/ /usr/local/bin/
COPY --from=build /workdir/frontend/dist frontend/dist

COPY config.toml config.toml

ENV OJ_LAB_SERVICE_ENV='production'
ENV DATABASE_DSN='user=postgres password=postgres host=host.docker.internal port=5432 dbname=oj_lab sslmode=disable TimeZone=Asia/Shanghai'
ENV REDIS_HOSTS='host.docker.internal:6379'
ENV MINIO_ENDPOINT='http://host.docker.internal:9000'
ENV SERVICE_MODE="release"

EXPOSE 8080
CMD [ "web_server" ]