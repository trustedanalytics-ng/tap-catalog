APDIR=$(shell go list ./... | grep -v /vendor/)
GOBIN=$(GOPATH)/bin

build:
	CGO_ENABLED=0 go install -tags netgo ${APDIR}
	go fmt $(APDIR)

run: build_anywhere
	./application/tap-catalog

run-local: build
	BROKER_LOG_LEVEL=DEBUG PORT=8083 CATALOG_USER=admin CATALOG_PASS=password ${GOPATH}/bin/tap-catalog

docker_build: build_anywhere
	docker build -t tap-catalog .

push_docker: docker_build
	docker tag -f tap-catalog $(REPOSITORY_URL)/tap-catalog:latest
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
	$(GOBIN)/govendor update github.com/trustedanalytics/...
	rm -Rf vendor/github.com/trustedanalytics/tap-catalog
	@echo "Done"

verify_gopath:
	@if [ -z "$(GOPATH)" ] || [ "$(GOPATH)" = "" ]; then\
		echo "GOPATH not set. You need to set GOPATH before run this command";\
		exit 1 ;\
	fi

tests: verify_gopath
	go test --cover $(APP_DIR_LIST)
	
prepare_dirs:
	mkdir -p ./temp/src/github.com/trustedanalytics/tap-catalog
	$(eval REPOFILES=$(shell pwd)/*)
	ln -sf $(REPOFILES) temp/src/github.com/trustedanalytics/tap-catalog

build_anywhere: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics/tap-catalog/... | grep -v /vendor/))
	GOPATH=$(GOPATH) CGO_ENABLED=0 go build -tags netgo $(APP_DIR_LIST)
	rm -Rf application && mkdir application
	cp -RL ./tap-catalog ./application/tap-catalog
	rm -Rf ./temp
