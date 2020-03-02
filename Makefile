.PHONY: deps test int-test unit-test mocks cover sync-coveralls

deps:
	go mod download

test: deps
	go test -cover `go list ./... | grep -v test`

integration: deps
	docker-compose -f ./test/integration/docker-compose.yml up --build --abort-on-container-exit
	docker-compose -f ./test/integration/docker-compose.yml down --volumes

covertest: deps
	go test  -coverprofile=coverage.out `go list ./... | grep -v test`
	go tool cover -html=coverage.out

sync-coveralls: deps
	go test  -coverprofile=coverage.out `go list ./... | grep -v test`
	goveralls -coverprofile=coverage.out -reponame=go-webdriver -repotoken=${COVERALLS_GO_WEBDRIVER_TOKEN} -service=local

mocks: deps
	mockgen -package=protocol -destination=pkg/protocol/client_mock.go -source=pkg/protocol/client.go
	mockgen -package=protocol -destination=pkg/protocol/session_mock.go -source=pkg/protocol/session.go
	mockgen -package=protocol -destination=pkg/protocol/timeouts_mock.go -source=pkg/protocol/timeouts.go
	mockgen -package=protocol -destination=pkg/protocol/navigation_mock.go -source=pkg/protocol/navigation.go
	mockgen -package=httpclient -destination=pkg/httpclient/client_mock.go -source=pkg/httpclient/client.go

