ARG GO_VERSION

FROM golang:${GO_VERSION}
ARG BUF_VERSION=1.31.0

WORKDIR /tmp/buf
RUN curl -sSL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(uname -s)-$(uname -m)" -o "/usr/local/bin/buf"
RUN chmod +x "/usr/local/bin/buf"

WORKDIR /usr/src/days
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go
RUN go install connectrpc.com/connect/cmd/protoc-gen-connect-go

WORKDIR ${GOPATH}

ENTRYPOINT ["/usr/local/bin/buf"]
