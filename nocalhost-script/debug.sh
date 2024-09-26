#! /bin/sh

if ! command -v dlv &> /dev/null; then
  go get -u github.com/go-delve/delve/cmd/dlv && go install github.com/go-delve/delve/cmd/dlv
else
  #dlv --headless --log --listen :9009 --api-version 2 --accept-multiclient debug sync_tencent_bill_to_db.go
  dlv --headless --log --listen :9009 --api-version 2 --accept-multiclient debug sync_aliyun_bill_to_db.go
fi
