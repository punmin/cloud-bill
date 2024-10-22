build:
	mkdir -p cloud-bill-bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloud-bill-bin/cloud-bill main.go
	cp config-tmpl.yaml cloud-bill-bin/config.yaml
	cp readme.md cloud-bill-bin/
	tar -zcvf cloud-bill-bin.tar.gz cloud-bill-bin
	rm -rf cloud-bill-bin

docker_build:
	docker build -t cloud-bill .

clean:
	rm -rf cloud-bill-bin	cloud-bill-bin.tar.gz