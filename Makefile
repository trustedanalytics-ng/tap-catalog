APDIR=$(shell go list ./... | grep -v /vendor/)
GOBIN=$(GOPATH)/bin

build:
	CGO_ENABLED=0 go install -tags netgo ${APDIR}
	go fmt $(APDIR)

run: build
	${GOPATH}/bin/tap-catalog

run-local: build
	CATALOG_PORT=8083 CATALOG_USER=admin CATALOG_PASS=admin ${GOPATH}/bin/tap-catalog

docker_build: build
	rm -Rf application && mkdir application
	cp -Rf $(GOBIN)/tap-catalog application/
	docker build -t tap-catalog .

push_docker: docker_build
	docker tag tap-catalog $(REPOSITORY_URL)/tap-catalog:latest
	docker push $(REPOSITORY_URL)/tap-catalog:latest

kubernetes_deploy:
	kubectl create -f configmap.yaml
	kubectl create -f service.yaml
	kubectl create -f deployment.yaml

deps_update: verify_gopath
	$(GOBIN)/govendor remove +all
	$(GOBIN)/govendor add +external
	@echo "Done"

bin/govendor: verify_gopath
	go get -v -u github.com/kardianos/govendor

deps_fetch_specific: bin/govendor
	@if [ "$(DEP_URL)" = "" ]; then\
		echo "DEP_URL not set. Run this comand as follow:";\
		echo " make deps_fetch_specific DEP_URL=github.com/nu7hatch/gouuid";\
	exit 1 ;\
	fi
	@echo "Fetchinf specific deps in newest versions"

	$(GOBIN)/govendor fetch -v $(DEP_URL)

verify_gopath:
	@if [ -z "$(GOPATH)" ] || [ "$(GOPATH)" = "" ]; then\
		echo "GOPATH not set. You need to set GOPATH before run this command";\
		exit 1 ;\
	fi

tests: verify_gopath
	go test --cover $(APP_DIR_LIST)