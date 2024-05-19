FROM golang:latest as build

COPY . /oj-lab-platform-build

WORKDIR /oj-lab-platform-build
RUN apt update && apt install -y make zip curl
RUN make build
RUN ./scripts/update-frontend-dist.sh /oj-lab-platform-build/frontend_dist


FROM ubuntu:latest

WORKDIR /workdir

COPY --from=build /oj-lab-platform-build/bin/web_server /usr/local/bin/web_server
COPY --from=build /oj-lab-platform-build/frontend_dist /workdir/frontend_dist

COPY workdirs/docker/config.toml /workdir/config.toml

ENV OJ_LAB_SERVICE_ENV=production
ENV OJ_LAB_WORKDIR=/workdir
EXPOSE 8080
CMD [ "web_server" ]