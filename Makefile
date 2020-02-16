
.PHONY: dep
dep: ## Install dependencies
	@go mod download
	@go mod vendor

.PHONY: test
test: dep ## Run test
	@go test -coverprofile=coverage.out ./... -v
	@go tool cover -html=coverage.out
