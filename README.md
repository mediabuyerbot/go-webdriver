# Go-WebDriver :: [W3C-Specification](https://w3c.github.io/webdriver/)
### Work in progress...

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

## Protocol implementation

| [**Session specification**](https://w3c.github.io/webdriver/#sessions)         | Chrome        | Firefox  |
| -----------------------------------------------------------------------------  | :------------:| :-------:|
| [New Session](https://w3c.github.io/webdriver/#new-session)                    |  &#10003;     | &#10003; |
| [Delete Session](https://w3c.github.io/webdriver/#delete-session)              |  &#10003;     | &#10003; |
| [Status](https://w3c.github.io/webdriver/#status)                              |  &#10003;     | &#10003; |
| | | |
| [**Timeouts specification**](https://w3c.github.io/webdriver/#timeouts)        | | |
| [Get Timeouts](https://w3c.github.io/webdriver/#get-timeouts)                  |  &#10003;     | &#10003; |
| [Set Timeouts](https://w3c.github.io/webdriver/#set-timeouts)                  |  &#10003;     | &#10003; |
| | | |

