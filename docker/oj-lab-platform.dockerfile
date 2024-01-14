FROM golang:latest as build

WORKDIR /oj-lab-platform-build

COPY go.mod /oj-lab-platform-build/go.mod
COPY go.sum /oj-lab-platform-build/go.sum
COPY scripts/ /oj-lab-platform-build/scripts/
COPY Makefile /oj-lab-platform-build/Makefile

COPY src/application/ /oj-lab-platform-build/src/application/
COPY src/core/ /oj-lab-platform-build/src/core/
COPY src/service/ /oj-lab-platform-build/src/service/

RUN apt update && apt install -y make zip curl
RUN make build
RUN make get-front


FROM ubuntu:latest

WORKDIR /workspace

COPY --from=build /oj-lab-platform-build/artifacts/bin/service /usr/local/bin/oj-lab-service
COPY --from=build /oj-lab-platform-build/artifacts/oj-lab-front /workspace/artifacts/oj-lab-front

COPY environment/configs/production.toml /workspace/environment/configs/production.toml

ENV OJ_LAB_SERVICE_ENV=production
ENV OJ_LAB_PROJECT_ROOT=workspace
EXPOSE 8080
CMD [ "oj-lab-service" ]