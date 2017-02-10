# Copyright (c) 2016 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
APP_DIR_LIST=$(shell go list ./... | grep -v /vendor/)
GOBIN=$(GOPATH)/bin
APP_NAME=tap-catalog

build: verify_gopath
	go fmt $(APP_DIR_LIST)
	CGO_ENABLED=0 go install -ldflags "-w" -tags netgo .
	mkdir -p application && cp -f $(GOBIN)/$(APP_NAME) ./application/$(APP_NAME)

run: build_anywhere
	./application/tap-catalog

run-local: build
	CORE_ORGANIZATION=default BROKER_LOG_LEVEL=DEBUG ETCD_CATALOG_HOST=localhost ETCD_CATALOG_PORT=2379 PORT=8084 CATALOG_USER=admin CATALOG_PASS=password ${GOPATH}/bin/tap-catalog

docker_build: build_anywhere
	docker build -t $(REPOSITORY_URL)/tap-catalog:latest .

push_docker: docker_build
	docker push $(REPOSITORY_URL)/tap-catalog:latest

kubernetes_deploy: docker_build
	kubectl create -f configmap.yaml
	kubectl create -f service.yaml
	kubectl create -f deployment.yaml

kubernetes_update: docker_build
	kubectl delete -f deployment.yaml
	kubectl create -f deployment.yaml

bin/govendor: verify_gopath
	go get -v -u github.com/kardianos/govendor

deps_fetch_specific: bin/govendor
	@if [ "$(DEP_URL)" = "" ]; then\
		echo "DEP_URL not set. Run this comand as follow:";\
		echo " make deps_fetch_specific DEP_URL=github.com/nu7hatch/gouuid";\
	exit 1 ;\
	fi
	@echo "Fetching specific dependency in newest versions"
	$(GOBIN)/govendor fetch -v $(DEP_URL)

deps_update_tap: verify_gopath
	$(GOBIN)/govendor update github.com/trustedanalytics-ng/...
	$(GOBIN)/govendor remove github.com/trustedanalytics-ng/$(APP_NAME)/...
	@echo "Done"

verify_gopath:
	@if [ -z "$(GOPATH)" ] || [ "$(GOPATH)" = "" ]; then\
		echo "GOPATH not set. You need to set GOPATH before run this command";\
		exit 1 ;\
	fi

test: verify_gopath
	CGO_ENABLED=0 go test --cover -ldflags "-w" -tags netgo $(APP_DIR_LIST)
	
prepare_dirs:
	mkdir -p ./temp/src/github.com/trustedanalytics-ng/tap-catalog
	$(eval REPOFILES=$(shell pwd)/*)
	ln -sf $(REPOFILES) temp/src/github.com/trustedanalytics-ng/tap-catalog

build_anywhere: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics-ng/tap-catalog/... | grep -v /vendor/))
	GOPATH=$(GOPATH) CGO_ENABLED=0 go build -tags netgo $(APP_DIR_LIST)
	rm -Rf application && mkdir application
	cp -RL ./tap-catalog ./application/tap-catalog
	rm -Rf ./temp

mock_update:
	$(GOBIN)/mockgen -source=data/data_repository.go -package=data -destination=data/data_repository_mock.go
	$(GOBIN)/mockgen -source=etcd/etcd_client.go -package=etcd -destination=etcd/etcd_client_mock.go
	$(GOBIN)/mockgen -source=vendor/github.com/coreos/etcd/client/keys.go -package=client -destination=vendor/github.com/coreos/etcd/client/keys_mock.go
	./add_license.sh
