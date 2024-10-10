#! /bin/sh

if [ ! -f "/go/bin/dlv" ]; then
  go get -u github.com/go-delve/delve/cmd/dlv && go install github.com/go-delve/delve/cmd/dlv
  dlv --headless --log --listen :9009 --api-version 2 --accept-multiclient debug main.go -- -c aws -m 2024-08
else
  dlv --headless --log --listen :9009 --api-version 2 --accept-multiclient debug main.go -- -c aws -m 2024-08
fi
