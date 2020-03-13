package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/mitchellh/go-ps"

	"github.com/mediabuyerbot/go-webdriver/pkg/chromedriver"
	bin "github.com/mediabuyerbot/go-webdriver/third_party/drivers"
)

var (
	port = flag.Int("port", 9515, "port to listen on")
)

func main() {
	flag.Parse()

	driver, err := chromedriver.New(bin.ChromeDriver64(),
		chromedriver.WithPort(*port),
		chromedriver.WithStderr(log.Writer()),
		// chromedriver.WithStdout(log.Writer()),
		chromedriver.WithRunHook(func(pid int) {
			fmt.Printf("ChromeDriver is running port=%d, portIsOpen=%v, proc=%d, procIsAlive=%v\n\n",
				*port, portIsOpen(*port), pid, procIsRunning(pid))
		}),
		chromedriver.WithStopHook(func(pid int) {
			fmt.Printf("ChromeDriver is stopped port=%d, proc=%d, procIsAlive=%v\n\n",
				*port, pid, procIsRunning(pid))
		}),
	)
	if err != nil {
		exitWithError(err)
	}

	ctx := context.Background()
	go func() {
		if err := driver.Run(ctx); err != nil {
			exitWithError(err)
		}
	}()

	terminate(func() error {
		if err := driver.Stop(ctx); err != nil {
			return err
		}
		return nil
	})
}

func exitWithError(err error) {
	fmt.Println("ERROR:", err)
	os.Exit(1)
}

func portIsOpen(port int) (status bool) {
	host := ":" + strconv.Itoa(port)
	conn, err := net.Listen("tcp", host)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func procIsRunning(pid int) bool {
	osp, err := ps.FindProcess(pid)
	if err != nil || osp == nil {
		return false
	}
	return osp.Pid() == pid
}

func terminate(fn func() error) {
	done := make(chan error)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- fn()
	}()
	err := <-done
	if err != nil {
		exitWithError(err)
	}
	time.Sleep(time.Second)
}
