deps:
	go get -d -v ./...

updatedeps:
	go get -d -v -u -f ./...

testdeps:
	go get -d -v -t ./...

updatetestdeps:
	go get -d -v -t -u -f ./...

build: deps
	go build ./...

install: deps
	go install ./...

lint: testdeps
	go get -v github.com/golang/lint/golint
	golint ./...

vet: testdeps
	go vet ./...

errcheck: testdeps
	go get -v github.com/kisielk/errcheck
	errcheck ./...

pretest: lint

test: testdeps pretest
	go test ./...

clean:
	go clean -i ./...

.PHONY: \
	deps \
	updatedeps \
	testdeps \
	updatetestdeps \
	build \
	install \
	lint \
	vet \
	errcheck \
	pretest \
	test \
	clean
