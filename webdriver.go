package webdriver

type Webdriver struct {
}

func (w *Webdriver) NewBrowser() (*Browser, error) {
	return &Browser{}, nil
}

func (w *Webdriver) Close() error {
	return nil
}
