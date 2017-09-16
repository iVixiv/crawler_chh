package home

import (
	"github.com/PuerkitoBio/goquery"
	"crawler_chh/browser"
	"strings"
	"github.com/op/go-logging"
	"crawler_chh/utils"
	"crawler_chh/cralwer"
)

const (
	HOME = "forum.php"
)

var (
	cBrowser = new(browser.Chh_browser)
	log      = logging.MustGetLogger("chh_home")
)

func (this Home) Crawler() (*[]Result, error) {
	doc, err := cBrowser.LoadHtmlDoc(cralwer.BASE_URL + HOME)
	if err != nil {
		return nil, err
	}
	return clean(doc), nil
}

func clean(doc *goquery.Document) (*[]Result) {
	result := make([]Result, 0)
	doc.Find(".bm.bmw.fl").Each(func(i int, selection *goquery.Selection) {
		cleanOne(selection, &result)
	})
	return &result
}

func cleanOne(selection *goquery.Selection, results *[]Result) {
	title := selection.Find(".subforumshow.bm_h.cl").Text()
	selection.Find("li").Each(func(i int, selection *goquery.Selection) {
		num := selection.Find(".num").Text()
		content := selection.Text()
		content = strings.Replace(content, num, "", -1)
		url, _ := selection.Find("a").Attr("href")
		r := Result{
			Num:        utils.FindUInt(num),
			ChildName:  utils.TrimDoc(content),
			Url:        utils.TrimDoc(url),
			ParentName: utils.TrimDoc(title),
		}
		*results = append(*results, r)
	})
}
