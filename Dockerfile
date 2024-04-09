FROM golang:latest as build

COPY go.mod /oj-lab-platform-build/go.mod
COPY go.sum /oj-lab-platform-build/go.sum
COPY scripts/ /oj-lab-platform-build/scripts/
COPY Makefile /oj-lab-platform-build/Makefile

COPY src/application/ /oj-lab-platform-build/src/application/
COPY src/core/ /oj-lab-platform-build/src/core/
COPY src/service/ /oj-lab-platform-build/src/service/

WORKDIR /oj-lab-platform-build
RUN apt update && apt install -y make zip curl
RUN make build
RUN ./scripts/update-frontend-dist.sh /oj-lab-platform-build/frontend_dist


FROM ubuntu:latest

WORKDIR /workdir

COPY --from=build /oj-lab-platform-build/bin/service /usr/local/bin/oj-lab-service
COPY --from=build /oj-lab-platform-build/frontend_dist /workdir/frontend_dist

COPY workdirs/docker/config.toml /workdir/config.toml

ENV OJ_LAB_SERVICE_ENV=production
ENV OJ_LAB_WORKDIR=/workdir
EXPOSE 8080
CMD [ "oj-lab-service" ]