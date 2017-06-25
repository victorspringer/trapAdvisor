test:
	export DB_USER="root" && \
	export DB_PASS="" && \
	export DB_ADDR="localhost" && \
	export DB_PORT="3306" && \
	go test -v -race $$(go list ./... | grep -v /vendor/)
test-fast:
	export DB_USER="root" && \
	export DB_PASS="" && \
	export DB_ADDR="localhost" && \
	export DB_PORT="3306" && \
	go test -v $$(go list ./... | grep -v /vendor/)
build:
	go get && go build -v -o bin/application
start:
	export ENV="production" && \
	export DOMAIN="localhost:8080" && \
	export PORT="8080" && \
	export DB_USER="root" && \
	export DB_PASS="" && \
	export DB_ADDR="localhost" && \
	export DB_PORT="3306" && \
	export FB_CLIENT_ID="" && \
	export FB_CLIENT_SECRET="" && \
	go get && go build -v -o bin/application && ./bin/application
start-dev:
	export ENV="development" && \
	export DOMAIN="localhost:8080" && \
	export PORT="8080" && \
	export DB_USER="root" && \
	export DB_PASS="" && \
	export DB_ADDR="localhost" && \
	export DB_PORT="3306" && \
	export FB_CLIENT_ID="" && \
	export FB_CLIENT_SECRET="" && \
	go get && go build -v -o bin/application && ./bin/application

