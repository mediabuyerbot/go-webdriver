.PHONY: dep test int-test unit-test mocks

dep:
	go mod vendor

unit-test: dep
	go test  -coverprofile=coverage.out `go list ./... | grep -v test` -v
	go tool cover -html=coverage.out

int-test: dep
	docker-compose -f ./test/docker-compose.yml up --build --abort-on-container-exit
	docker-compose -f ./test/docker-compose.yml down --volumes

mocks: dep
	mockgen -package=protocol -destination=pkg/protocol/client_mock.go -source=pkg/protocol/client.go
	mockgen -package=protocol -destination=pkg/protocol/session_mock.go -source=pkg/protocol/session.go
	mockgen -package=protocol -destination=pkg/protocol/timeouts_mock.go -source=pkg/protocol/timeouts.go
	mockgen -package=httpclient -destination=pkg/httpclient/client_mock.go -source=pkg/httpclient/client.go

