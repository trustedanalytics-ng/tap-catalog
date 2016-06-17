export APDIR=$(go list ./... | grep -v /vendor/)

build:
	CGO_ENABLED=0 go install -tags netgo ${APDIR}

run: build
	${GOPATH}/bin/tap-catalog

run-local: build
	PORT=8181 CATALOG_USER=admin CATALOG_PASS=admin ${GOPATH}/bin/tap-catalog
