test:
	go test -v -race $$(go list ./... | grep -v /vendor/)
test-fast:
	go test -v $$(go list ./... | grep -v /vendor/)
build:
	go get && go build -v -o dist/trapAdvisor
start:
	go get && go build -v -o dist/trapAdvisor && ./dist/trapAdvisor