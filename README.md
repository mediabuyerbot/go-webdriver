# Go-WebDriver :: [W3C-Specification](https://w3c.github.io/webdriver/)

## Work in progress...

#### Installation
```ssh
go get github.com/mediabuyerbot/go-webdriver
```

#### Commands
```ssh
# Get golang dependencies 
$ make dep  

# Running unit tests 
$ make unit-test

# Running integration tests
$ make int-test

# Generate mocks
$ make mocks
```

## WebDriver protocol implementation
|                                                                                | Chrome        | Firefox  |
| -----------------------------------------------------------------------------  | :------------:| :-------:|
| [New Session](https://w3c.github.io/webdriver/#new-session)                    |  &#10003;     | &#10003; |
| [Delete Session](https://w3c.github.io/webdriver/#delete-session)              |  &#10003;     | &#10003; |
| [Status](https://w3c.github.io/webdriver/#status)                              |  &#10003;     | &#10003; |
| [Get Timeouts](https://w3c.github.io/webdriver/#get-timeouts)                  |  &#10003;     | &#10003; |
| [Set Timeouts](https://w3c.github.io/webdriver/#set-timeouts)                  |  &#10003;     | &#10003; |

