GOLINT=`go list -f {{.Target}} golang.org/x/lint/golint`

help:  ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

format: ## Run GOFMT to format code
	echo "Running GOFMT"
	go fmt ./...

lint: ## Run GOVET and GOLINT to check code quality
	echo "Running GOVET and GOLINT"
	go vet ./... && \
	$(GOLINT) ./...

test: ## Run tests
	echo "Running application test"
	go test ./...
	
build: lint test ## Build application
	echo "Running build"
	go build -o flightrouter
