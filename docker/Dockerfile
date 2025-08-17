FROM golang:1.24-alpine

RUN apk add --update \
    curl \
    libc6-compat

RUN curl -L https://github.com/Yandex-Practicum/go-autotests/releases/latest/download/metricstest -O --output-dir ${GOPATH}/bin
RUN chmod +x ${GOPATH}/bin/metricstest

RUN curl -L https://github.com/Yandex-Practicum/go-autotests/releases/latest/download/statictest -O --output-dir ${GOPATH}/bin
RUN chmod +x ${GOPATH}/bin/statictest

RUN curl -L https://github.com/Yandex-Practicum/go-autotests/releases/latest/download/random -O --output-dir ${GOPATH}/bin
RUN chmod +x ${GOPATH}/bin/random

RUN curl -L https://github.com/ktigay/multichecker/releases/latest/download/multichecker -O --output-dir ${GOPATH}/bin
RUN chmod +x ${GOPATH}/bin/multichecker

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ${GOPATH}/bin v2.1.6

RUN go install mvdan.cc/gofumpt@latest

RUN go install github.com/golang/mock/mockgen@latest

WORKDIR /app

COPY . .
RUN go mod download
