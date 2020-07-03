package webdriver

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"

	bin "github.com/mediabuyerbot/go-webdriver/third_party/drivers"

	"github.com/mediabuyerbot/go-webdriver/pkg/chromedriver"
)

const (
	EnvChromeDriverPath = "CHROMEDRIVER_PATH"
)

func Chrome(opts *ChromeOptionsBuilder) (*Browser, error) {
	if opts == nil {
		opts = ChromeOptions()
	}
	port, err := freePort()
	if err != nil {
		return nil, err
	}
	done := make(chan error)
	logPath, err := homedir.Expand("~/chrome-driver.log")
	if err != nil {
		return nil, err
	}
	driverPath := os.Getenv(EnvChromeDriverPath)
	if len(driverPath) == 0 {
		driverPath = bin.ChromeDriver64()
	}
	driver, err := chromedriver.New(driverPath,
		chromedriver.WithPort(port),
		chromedriver.WithStderr(log.Writer()),
		chromedriver.WithWhitelistedIps([]string{"0.0.0.0"}),
		chromedriver.WithLogPath(logPath),
		chromedriver.WithLogLevel(chromedriver.Severe),
		chromedriver.WithRunHook(func(pid int) {
			done <- nil
		}),
	)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	go func() {
		if err := driver.Run(ctx); err != nil {
			select {
			case done <- err:
				log.Printf("[ERROR] chromedriver failed to start\n%v\n", err)
			default:
				log.Printf("[ERROR] chromedriver \n%v\n", err)
			}
		}
		close(done)
	}()
	err = <-done
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("http://localhost:%d", port)
	sess, err := NewSession(ctx, addr, opts.Build())
	if err != nil {
		_ = driver.Stop(ctx)
		return nil, err
	}
	return &Browser{
		ctx:    ctx,
		driver: driver,
		sess:   sess,
	}, nil
}
