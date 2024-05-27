ARG GO_VERSION

# ===== build go binary =====
FROM golang:${GO_VERSION}-alpine3.19 as go-builder

WORKDIR /go/src/github.com/karamaru-alpha/days

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go install github.com/daichirata/hammer

COPY cmd/db-schema-sync cmd/db-schema-sync

RUN CGO_ENABLED=0 go build -o db-schema-sync -trimpath -ldflags '-s -w' cmd/db-schema-sync/main.go

# ==== build docker image ====
FROM alpine:3.20.0

WORKDIR /usr/src/days

COPY --from=go-builder /go/src/github.com/karamaru-alpha/days/db-schema-sync db-schema-sync
COPY --from=go-builder /go/bin/hammer /usr/local/bin/hammer
COPY /db/ddl db/ddl

ENTRYPOINT ["/usr/src/days/db-schema-sync"]
