# This Dockerfile separated to multiple layers in purpose to increase speed.
# As each layer in Docker is cached, in most cases it need to rebuild only `build` layer.
# It does not depend on CI/CD or developer's machine enviroment, so this allows
# to build tested production-ready image anywhere anytime. If the development process includes 
# complicated CI/CD logic, it should be divided to Dockerfile.build and Dockerfile.ci.

# base layer prepares most common, code independent image
FROM golang:1.12-alpine3.9 AS base

ENV CGO_ENABLED=0 
ENV GOOS=linux
ENV GOARCH=amd64 

RUN apk add git
WORKDIR /go/src/github.com/kfreiman/imageservice
RUN go get -u golang.org/x/lint/golint

# deps layer downloads vendor dependencies
FROM base AS deps 

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
RUN go mod download

# build layer. Build binary, test, lint, vet and other. 
FROM deps AS build

COPY . .

RUN go generate ./...
RUN go vet ./...
RUN golint -set_exit_status $(go list ./...)
RUN out=$(go fmt ./...) && if [[ -n "$out" ]]; then echo "$out"; exit 1; fi
RUN go test ./...
RUN go build -v -o /app/server github.com/kfreiman/imageservice/cmd/server

# produce final minimalistic alpine based image
FROM alpine:3.9
WORKDIR /app
COPY --from=build /app/server /app/server
CMD ["/app/server"]
