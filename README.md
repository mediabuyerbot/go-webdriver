# go-webdriver  [![Coverage Status](https://coveralls.io/repos/github/mediabuyerbot/go-webdriver/badge.svg?branch=master&t5)](https://coveralls.io/github/mediabuyerbot/go-webdriver?branch=master)

## Work in progress...

## Contents
- [Installation](#installation)
- [Commands](#commands)
  + [Build dependencies](#build-dependencies)
  + [Run test](#run-test)
  + [Run test with coverage profile](#run-test-with-coverage-profile)
  + [Run integration test](#run-integration-test)
  + [Run sync coveralls](#run-sync-coveralls)
  + [Build mocks](#build-mocks) 
  + [Download ChromeDriver, GeckoDriver](#download-chromedriver-geckodriver-third_partydrivers)
- [ChromeOptions docs](#chromeoptions-docs)
- [FirefoxOptions docs](#firefoxoptions-docs)
- [Protocol implementation](#protocol-implementation)  
  + [Session](#session)
  + [Timeouts](#timeouts)
  + [Navigation](#navigation)
  + [Context](#context)
  + [Cookies](#cookies)
  + [Document](#document)
  + [Screen capture](#screen-capture)
  + [User prompts](#user-prompts)
  + [Elements](#elements)

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
#### Download ChromeDriver, GeckoDriver third_party/drivers
```shell script
make download-drivers
```

## ChromeOptions docs 
+ [V8 dev](https://v8.dev/)
+ [Chromium command line switches](https://peter.sh/experiments/chromium-command-line-switches/)
+ [Chromium command line switches](https://chromium.googlesource.com/chromium/src/+/master/chrome/common/chrome_switches.cc)
+ [Chromium command line prefs name](https://chromium.googlesource.com/chromium/src/+/master/chrome/common/pref_names.cc)

## FirefoxOptions docs 
TODO

## Protocol implementation
### Session
| Specification                                                                 | Example | Chrome        | Firefox  |
| ----------------------------------------------------------------------------- |-------- | :------------:| :-------:|
| [New Session](https://w3c.github.io/webdriver/#new-session)                   |         |  &#10003;     | &#10003; |
| [Delete Session](https://w3c.github.io/webdriver/#delete-session)             |         |  &#10003;     | &#10003; |
| [Status](https://w3c.github.io/webdriver/#status)                             |         |  &#10003;     | &#10003; |

### Timeouts 
| Specification                                                                 | Example | Chrome        | Firefox  |
| ----------------------------------------------------------------------------- |-------- | :------------:| :-------:|
| [Get Timeouts](https://w3c.github.io/webdriver/#get-timeouts)                 |         |  &#10003;     | &#10003; |
| [Set Timeouts](https://w3c.github.io/webdriver/#set-timeouts)                 |         |  &#10003;     | &#10003; |

### Navigation
|  Specification                                                                 | Example       | Chrome        | Firefox  |
| -----------------------------------------------------------------------------  | ------------- | :------------:| :-------:|
| [Navigate To](https://w3c.github.io/webdriver/#navigate-to)                    |               |  &#10003;     | &#10003; |
| [Get Current URL](https://w3c.github.io/webdriver/#get-current-url)            |               |  &#10003;     | &#10003; |
| [Back](https://w3c.github.io/webdriver/#back)                                  |               |  &#10003;     | &#10003; |
| [Forward](https://w3c.github.io/webdriver/#forward)                            |               |  &#10003;     | &#10003; |
| [Refresh](https://w3c.github.io/webdriver/#refresh)                            |               |  &#10003;     | &#10003; |
| [Get Title](https://w3c.github.io/webdriver/#get-title)                        |               |  &#10003;     | &#10003; |

### Context
| Specification                                                                  | Example       | Chrome        | Firefox  |
| -----------------------------------------------------------------------------  | ------------- | :------------:| :-------:|
| [Get Window Handle](https://w3c.github.io/webdriver/#get-window-handle)        |               |  &#10003;     | &#10003; |
| [Close Window](https://w3c.github.io/webdriver/#close-window)                  |               |  &#10003;     | &#10003; |
| [Switch To Window](https://w3c.github.io/webdriver/#switch-to-window)          |               |  &#10003;     | &#10003; |
| [Get Window Handles](https://w3c.github.io/webdriver/#get-window-handles)      |               |  &#10003;     | &#10003; |
| [New Window](https://w3c.github.io/webdriver/#new-window)                      |               |  &#10003;     | &#10003; |
| [Switch To Frame](https://w3c.github.io/webdriver/#switch-to-frame)            |               |  &#10003;     | &#10003; |
| [Switch To Parent Frame](https://w3c.github.io/webdriver/#switch-to-parent-frame)|             |  &#10003;     | &#10003; |
| [Get Window Rect](https://w3c.github.io/webdriver/#get-window-rect)            |               |  &#10003;     | &#10003; |
| [Set Window Rect](https://w3c.github.io/webdriver/#set-window-rect)            |               |  &#10003;     | &#10003; |
| [Maximize Window](https://w3c.github.io/webdriver/#maximize-window)            |               |  &#10003;     | &#10003; |
| [Minimize Window](https://w3c.github.io/webdriver/#minimize-window)            |               |  &#10003;     | &#10003; |
| [Fullscreen Window](https://w3c.github.io/webdriver/#fullscreen-window)        |               |  &#10003;     | &#10003; |

### Cookies
| Specification                                                                  | Example       | Chrome        | Firefox  |
| -----------------------------------------------------------------------------  | ------------- | :------------:| :-------:|
| [Get All Cookies](https://w3c.github.io/webdriver/#get-all-cookies)            |               |  &#10003;     | &#10003; |
| [Get Named Cookie](https://w3c.github.io/webdriver/#get-named-cookie)          |               |  &#10003;     | &#10003; |
| [Add Cookie](https://w3c.github.io/webdriver/#add-cookie)                      |               |  &#10003;     | &#10003; |
| [Delete Cookie](https://w3c.github.io/webdriver/#delete-cookie)                |               |  &#10003;     | &#10003; |
| [Delete All Cookies](https://w3c.github.io/webdriver/#delete-all-cookies)      |               |  &#10003;     | &#10003; |

### Document
| Specification                                                                  | Example       | Chrome        | Firefox  |
| -----------------------------------------------------------------------------  | ------------- | :------------:| :-------:|
| [Get Page Source](https://w3c.github.io/webdriver/#get-page-source)            |               |  &#10003;     | &#10003; |
| [Execute Script](https://w3c.github.io/webdriver/#execute-script)              |               |  &#10003;     | &#10003; |
| [Execute Async Script](https://w3c.github.io/webdriver/#execute-async-script)  |               |  &#10003;     | &#10003; |

### Screen capture 
| Specification                                                                  | Example       | Chrome        | Firefox  |
| -----------------------------------------------------------------------------  | ------------- | :------------:| :-------:|
| [Take Screenshot](https://w3c.github.io/webdriver/#take-screenshot)            |               |  &#10003;     | &#10003; |
| [Take Element Screenshot](https://w3c.github.io/webdriver/#take-element-screenshot) |          |  &#10003;     | &#10003; |

### User prompts
| Specification                                                                  | Example       | Chrome        | Firefox  |
| -----------------------------------------------------------------------------  | ------------- | :------------:| :-------:|
| [Dismiss Alert](https://w3c.github.io/webdriver/#dismiss-alert)                |               |  &#10003;     | &#10003; |
| [Accept Alert](https://w3c.github.io/webdriver/#accept-alert)                  |               |  &#10003;     | &#10003; |
| [Get Alert Text](https://w3c.github.io/webdriver/#get-alert-text)              |               |  &#10003;     | &#10003; |
| [Send Alert Text](https://w3c.github.io/webdriver/#send-alert-text)            |               |  &#10003;     | &#10003; |

### Elements
| Specification                                                                  | Example       | Chrome        | Firefox  |
| -----------------------------------------------------------------------------  | ------------- | :------------:| :-------:|
| [Find Element](https://w3c.github.io/webdriver/#find-element)                  |               |  &#10003;     | &#10003; |
| [Find Elements](https://w3c.github.io/webdriver/#find-elements)                |               |  &#10003;     | &#10003; |
| [Find Element From Element](https://w3c.github.io/webdriver/#find-element-from-element)              |               |  &#10003;     | &#10003; |
| [Find Elements From Element](https://w3c.github.io/webdriver/#find-elements-from-element)            |               |  &#10003;     | &#10003; |
| [Get Active Element](https://w3c.github.io/webdriver/#get-active-element)           |               |  &#10003;     | &#10003; |
| [Is Element Selected](https://w3c.github.io/webdriver/#is-element-selected)         |               |  &#10003;     | &#10003; |
| [Get Element Attribute](https://w3c.github.io/webdriver/#get-element-attribute)     |               |  &#10003;     | &#10003; |
| [Get Element Property](https://w3c.github.io/webdriver/#get-element-property)       |               |  &#10003;     | &#10003; |
| [Get Element CSS Value](https://w3c.github.io/webdriver/#get-element-css-value)     |               |  &#10003;     | &#10003; |
| [Get Element Text](https://w3c.github.io/webdriver/#get-element-text)               |               |  &#10003;     | &#10003; |
| [Get Element Tag Name](https://w3c.github.io/webdriver/#get-element-tag-name)       |               |  &#10003;     | &#10003; |
| [Get Element Rect](https://w3c.github.io/webdriver/#get-element-rect)               |               |  &#10003;     | &#10003; |
| [Element Enabled](https://w3c.github.io/webdriver/#is-element-enabled)            |               |  &#10003;     | &#10003; |
| [Element Click](https://w3c.github.io/webdriver/#element-click)                     |               |  &#10003;     | &#10003; |
| [Element Clear](https://w3c.github.io/webdriver/#element-clear)                     |               |  &#10003;     | &#10003; |
| [Element Send Keys](https://w3c.github.io/webdriver/#element-send-keys)             |               |  &#10003;     | &#10003; |
