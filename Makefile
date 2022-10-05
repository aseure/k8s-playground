BIN := server
SRC_DIR := src
BUILD_DIR := dist
K3D_CONFIG := k3d/config.yaml

.PHONY: build-binaries
build-binaries:
	rm -rf ${BUILD_DIR}
	mkdir -p ${BUILD_DIR}
	GOOS=linux  GOARCH=amd64 go build -o  ./${BUILD_DIR}/linux-amd64/${BIN} ./${SRC_DIR}/...
	GOOS=darwin GOARCH=arm64 go build -o ./${BUILD_DIR}/darwin-arm64/${BIN} ./${SRC_DIR}/...

.PHONY: build-docker
build-docker: build-binaries
	docker build -t aseure/k8s-playground:latest .

.PHONY: publish-docker
publish-docker: build-docker
	docker push aseure/k8s-playground:latest

.PHONY: create-k8s-cluster
create-k8s-cluster:
	k3d cluster create --config ${K3D_CONFIG}

.PHONY: stop-k8s-cluster
stop-k8s-cluster:
	k3d cluster stop --config ${K3D_CONFIG}

.PHONY: delete-k8s-cluster
delete-k8s-cluster:
	k3d cluster delete --config ${K3D_CONFIG}
