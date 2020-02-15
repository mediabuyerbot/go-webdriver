package webdriver

type Browser struct {
	sessionID string
}

func (b *Browser) Cookie() *Cookie {
	return &Cookie{}
}
