FROM golang:1.22.4 AS builder

ARG OPERATOR_VERSION

RUN apk update && apk upgrade && \
    apk add --no-cache bash build-base

WORKDIR /opt/db-operator

# to reduce docker build time download dependency first before building
COPY go.mod .
COPY go.sum .
RUN go mod download

# build
COPY . .

RUN CGO_ENABLED=0 go build -tags build -o /usr/local/bin/db-operator -ldflags="-X 'github.com/db-operator/db-operator/internal/helpers/common.OperatorVersion=${OPERATOR_VERSION}'" cmd/main.go

FROM alpine:3.20
LABEL org.opencontainers.image.authors="Nikolai Rodionov<allanger@zohomail.com>"

ENV USER_UID=1001
ENV USER_NAME=db-operator

# # install operator binary
COPY --from=builder /usr/local/bin/db-operator /usr/local/bin/db-operator
COPY ./build/bin /usr/local/bin
RUN /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]
USER $USER_UID
