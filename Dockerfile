# Build stage
FROM golang:1.14-alpine as build

RUN apk update && apk upgrade \
  && apk add git ca-certificates \
  && rm -rf /var/cache/apk/*

WORKDIR /
ARG PORT
COPY . .
RUN set -ex; \
    SVC_VERSION=$(git --no-pager log -1 --pretty=format:%h); \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -o kubetimer -ldflags "-X main.version=${SVC_VERSION}" .

# Final image with no source code
FROM alpine:3.12

RUN apk update && apk upgrade \
  && apk add ca-certificates \
  && rm -rf /var/cache/apk/*

WORKDIR /
COPY --from=build /kubetimer .
EXPOSE ${PORT}
ENTRYPOINT /kubetimer
