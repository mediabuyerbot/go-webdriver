# go-webdriver  [![Coverage Status](https://coveralls.io/repos/github/mediabuyerbot/go-webdriver/badge.svg?branch=master)](https://coveralls.io/github/mediabuyerbot/go-webdriver?branch=master)

## Work in progress...

## Table of contents
- [Installation](#installation)
- [Commands](#commands)
  + [Build dependencies](#build-dependencies)
  + [Run test](#run-test)
  + [Run test with coverage profile](#run-test-with-coverage-profile)
  + [Run integration test](#run-integration-test)
  + [Run sync coveralls](#run-sync-coveralls)
  + [Build mocks](#build-mocks) 
- [Protocol implementation](#protocol-implementation)  

### Installation
```ssh
go get github.com/mediabuyerbot/go-webdriver
```

### Commands
#### Build dependencies
```shell script
make deps
```
#### Run test
```shell script
make test
```
#### Run test with coverage profile
```shell script
make covertest
```
#### Run integration test
```shell script
make integration
```
#### Run sync coveralls
```shell script
COVERALLS_GO_WEBDRIVER_TOKEN=${COVERALLS_REPO_TOKEN}
make sync-coveralls
```
#### Build mocks
```shell script
make mocks
```

## Protocol implementation

| Session specification                                                          | Chrome        | Firefox  | 
| -----------------------------------------------------------------------------  | :------------:| :-------:|
| [New Session](https://w3c.github.io/webdriver/#new-session)                    |  &#10003;     | &#10003; |
| [Delete Session](https://w3c.github.io/webdriver/#delete-session)              |  &#10003;     | &#10003; |
| [Status](https://w3c.github.io/webdriver/#status)                              |  &#10003;     | &#10003; |

| Timeouts specification                                                         | Chrome        | Firefox  |
| -----------------------------------------------------------------------------  | :------------:| :-------:|
| [Get Timeouts](https://w3c.github.io/webdriver/#get-timeouts)                  |  &#10003;     | &#10003; |
| [Set Timeouts](https://w3c.github.io/webdriver/#set-timeouts)                  |  &#10003;     | &#10003; |

| Navigation specification                                                       | Chrome        | Firefox  |
| -----------------------------------------------------------------------------  | :------------:| :-------:|
| [Navigate To](https://w3c.github.io/webdriver/#navigate-to)                    |  &#10003;     | &#10003; |
| [Get Current URL](https://w3c.github.io/webdriver/#get-current-url)            |  &#10003;     | &#10003; |
| [Back](https://w3c.github.io/webdriver/#back)                                  |  &#10003;     | &#10003; |
| [Forward](https://w3c.github.io/webdriver/#forward)                            |  &#10003;     | &#10003; |
| [Refresh](https://w3c.github.io/webdriver/#refresh)                            |  &#10003;     | &#10003; |
| [Get Title](https://w3c.github.io/webdriver/#get-title)                        |  &#10003;     | &#10003; |

| Context specification                                                          | Chrome        | Firefox  |
| -----------------------------------------------------------------------------  | :------------:| :-------:|
| [Get Timeouts](https://w3c.github.io/webdriver/#get-timeouts)                  |  &#10003;     | &#10003; |