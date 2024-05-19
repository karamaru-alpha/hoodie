#!/bin/sh

cd proto || exit $?

# buf lint
buf lint || exit $?

# store start time
touch -t $(date +%Y%m%d%H%M.%S) /tmp/buf-start-timestamp

# buf generate
call_buf_generate() {
  buf generate "$@" || exit $?
}
call_buf_generate --path options --template options.buf.gen.yaml
call_buf_generate --path enums --template enum.buf.gen.yaml
call_buf_generate --path entity/transaction --template transaction.buf.gen.yaml
call_buf_generate --path rpc/api --template api.buf.gen.yaml

# buf format
buf format -w

# delete old files
cd - || exit $?
find pkg/ \( -name '*.pb.go' -o -name '*.connect.go' -o -name '*.gen.go' \) ! -newer /tmp/buf-start-timestamp -exec rm {} \;
