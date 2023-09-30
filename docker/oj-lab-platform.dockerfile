FROM golang:latest as build

COPY application/ /usr/src/application/
COPY core/ /usr/src/core/
COPY service/ /usr/src/service/
COPY go.mod /usr/src/go.mod
COPY go.sum /usr/src/go.sum
COPY Makefile /usr/src/Makefile

WORKDIR /usr/src

RUN apt update && apt install -y make
RUN make build


FROM ubuntu:latest

COPY --from=build /usr/src/bin/service /usr/local/bin/oj-lab-service

RUN apt update && apt install -y make zip curl

RUN mkdir /workspace
COPY config/production.toml /workspace/config/production.toml
WORKDIR /workspace

COPY Makefile /workspace/Makefile
COPY script/ /workspace/script/
RUN make get-front
RUN apt install -y strace

ENV OJ_LAB_SERVICE_ENV=production
ENV OJ_LAB_PROJECT_ROOT=workspace
EXPOSE 8080
CMD [ "oj-lab-service" ]