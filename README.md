# go-webdriver

# design
``` 
browser, err := webdriver.NewBrowser("127.0.0.1:9090", 
   webdriver.DefaultDesiredLinuxCapabilities,
   webdriver.DefaultRequiredLinuxCapabilities
)
if err != nil {
   panic(err)
}
defer browser.Close()

browser.Cookies().Set(webdriver.Cookie{Name: "", Value: ""})

if err := browser.Navigation().To("http://google.com"); err != nil {
   panic(err)
}
```