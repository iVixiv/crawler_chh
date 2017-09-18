package browser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
	"strings"
	"fmt"
)

type Chh_browser struct {
}

func (chh_browser Chh_browser) LoadHtmlDoc(url string) (*goquery.Document, error) {
	driver := getPhoneWebDriver()
	if err := driver.Start(); err != nil {
		return nil, fmt.Errorf("driver.Start() : %s", err.Error())
	}
	defer driver.Stop()
	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		return nil, fmt.Errorf("driver.NewPage() : %s", err.Error())
	}
	if err := page.Navigate(url); err != nil {
		return nil, fmt.Errorf("driver.Navigate() : %s", err.Error())
	}
	//page.Screenshot("chh.png")
	html, err := page.HTML()
	if err != nil {
		return nil, fmt.Errorf("page.HTML() : %s", err.Error())
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	return doc, fmt.Errorf("goquery.NewDocumentFromReader() : %s", err.Error())
}
