.PHONY: dep test int-test unit-test mocks

MOCK_POSTFIX = mock

dep:
	go mod vendor

unit-test: dep
	go test  -coverprofile=coverage.out `go list ./... | grep -v test` -v
	go tool cover -html=coverage.out

int-test: dep
	docker-compose -f ./test/docker-compose.yml up --build --abort-on-container-exit
	docker-compose -f ./test/docker-compose.yml down --volumes

mocks: dep
	mockgen -package=protocol -destination=pkg/protocol/client_${MOCK_POSTFIX}.go -source=pkg/protocol/client.go
	mockgen -package=httpclient -destination=pkg/httpclient/client_${MOCK_POSTFIX}.go -source=pkg/httpclient/client.go

