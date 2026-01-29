.PHONY: help build docker_build docker_build_and_push clean

IMAGE ?= cloud-bill

help:
	@echo "使用以下命令进行构建镜像并推送"
	@echo "make docker_build_and_push IMAGE=your_registry/repo/cloud-bill:v1.1.1"

build:
	mkdir -p cloud-bill-bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloud-bill-bin/cloud-bill main.go
	cp config-tmpl.yaml cloud-bill-bin/config.yaml
	cp readme.md cloud-bill-bin/
	tar -zcvf cloud-bill-bin.tar.gz cloud-bill-bin
	rm -rf cloud-bill-bin

docker_build:
	docker build -t $(IMAGE) .

docker_build_and_push:
	docker build -t $(IMAGE) .
	docker push $(IMAGE)

clean:
	rm -rf cloud-bill-bin	cloud-bill-bin.tar.gz