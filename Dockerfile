FROM golang:1.23.0-alpine

# Enable go modules
ENV GO111MODULE=on

# Install dependencies
RUN apk add curl git build-base

# Install linter
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $HOME/bin v1.44.0

# Copy go mod files first and install dependencies to cache this layer
ADD ./go.mod ./go.sum /go/src/terraform-provider-opengpg/
WORKDIR /go/src/terraform-provider-opengpg
RUN go mod download

# Add source code
ADD . /go/src/terraform-provider-opengpg

# Build, test and lint
RUN go build -v && \
    go test -v ./... && \
    $HOME/bin/golangci-lint run
