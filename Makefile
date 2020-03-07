.PHONY: deps test  mocks cover sync-coveralls download-drivers

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

download-drivers:
	@cd ./third_party/drivers && ./download.sh

mocks: deps
	mockgen -package=protocol -destination=pkg/protocol/transport_mock.go -source=pkg/protocol/transport.go
	mockgen -package=protocol -destination=pkg/protocol/session_mock.go -source=pkg/protocol/session.go




