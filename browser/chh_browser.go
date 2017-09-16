package browser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
	"strings"
)

type Chh_browser struct {
}

func (chh_browser Chh_browser) LoadHtmlDoc(url string) (*goquery.Document, error) {
	driver := getPhoneWebDriver()
	if err := driver.Start(); err != nil {
		return nil, err
	}
	defer driver.Stop()
	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		return nil, err
	}
	if err := page.Navigate(url); err != nil {
		return nil, err
	}
	//page.Screenshot("chh.png")
	html, err := page.HTML()
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(strings.NewReader(html))
}
