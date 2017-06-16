test:
	go test -v -race $$(go list ./... | grep -v /vendor/)
test-fast:
	go test -v $$(go list ./... | grep -v /vendor/)
build:
	go get && go build -v -o bin/trapAdvisor
start:
	go get && go build -v -o bin/trapAdvisor && ./bin/trapAdvisor