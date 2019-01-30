ARCH=amd64
OS=linux
IMAGE=mx3d/netapp_exporter
VERSION=v0.1

build:
	GOOS=$(OS) GOARCH=$(ARCH) go build
	docker build -t $(IMAGE):$(VERSION) . 

push:
	docker push $(IMAGE):$(VERSION)
