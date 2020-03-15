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
	mockgen -package=w3c -destination=pkg/w3c/transport_mock.go -source=pkg/w3c/transport.go
	mockgen -package=w3c -destination=pkg/w3c/session_mock.go -source=pkg/w3c/session.go
	mockgen -package=w3c -destination=pkg/w3c/browser_options_mock.go -source=pkg/w3c/browser_options.go




