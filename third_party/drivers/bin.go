package bin

import (
	"path/filepath"
	"runtime"
)

func ChromeDriver64() (binpath string) {
	return lookPath(
		"chromedriver_mac64",
		"chromedriver_linux64",
	)
}

func GeckoDriver64() (binpath string) {
	return lookPath(
		"geckodriver_macos",
		"geckodriver_linux64",
	)
}

func lookPath(darwin string, linux string) (binpath string) {
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
