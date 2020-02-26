.PHONY: dep test test-int

dep: # install dependencies
	go mod tidy
	go mod download
	go mod vendor

unit-test: dep # run unit tests
	go test -coverprofile=coverage.out ./... -v
	go tool cover -html=coverage.out

int-test: dep #  run integration tests
	docker-compose -f ./test/docker-compose.yml up --build --abort-on-container-exit
	docker-compose -f ./test/docker-compose.yml down --volumes


