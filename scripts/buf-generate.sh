#!/bin/sh

cd proto || exit $?

# buf lint
buf lint || exit $?

# buf generate
call_buf_generate() {
  buf generate "$@" || exit $?
}
call_buf_generate --path options --template options.buf.gen.yaml
call_buf_generate --path enums --template enum.buf.gen.yaml
call_buf_generate --path rpc --template rpc.buf.gen.yaml
call_buf_generate --path entity/transaction --template transaction.buf.gen.yaml
