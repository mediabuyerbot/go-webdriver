package bin

import (
	"path/filepath"
	"runtime"
)

func ChromeDriver64() (binpath string) {
	return lookDriverPath(
		"chromedriver_mac64",
		"chromedriver_linux64",
	)
}

func GeckoDriver64() (binpath string) {
	return lookDriverPath(
		"geckodriver_macos",
		"geckodriver_linux64",
	)
}

func lookDriverPath(darwin string, linux string) (binpath string) {
	var (
		_, filename, _, _ = runtime.Caller(0)
		pkg               = filepath.Dir(filename)
	)
	switch runtime.GOOS {
	case "darwin":
		binpath = filepath.Join(pkg, darwin)
	default:
		binpath = filepath.Join(pkg, linux)
	}
	return binpath
}
