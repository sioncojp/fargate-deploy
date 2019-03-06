# build stage
FROM golang:latest AS build-env

ENV GO111MODULE auto

ADD . /src
WORKDIR /src
RUN make build

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/cmd/fargate-deploy /app/
ENTRYPOINT ./fargate-deploy
