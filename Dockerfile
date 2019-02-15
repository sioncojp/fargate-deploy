# build stage
FROM golang:alpine AS build-env

ENV GO111MODULE on

ADD . /src
RUN cd /src && go build -o fargate-deploy 

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/fargate-deploy /app/
ENTRYPOINT ./fargate-deploy
