package mod

import (
	"crawler_chh/browser"
	"strings"
	"crawler_chh/cralwer"
	"github.com/PuerkitoBio/goquery"
	"crawler_chh/utils"
	"github.com/op/go-logging"
)

var (
	cBrowser = new(browser.Chh_browser)
	log      = logging.MustGetLogger("chh_mod")
)

func (this Mod) Crawler(url string) (*Result, error) {
	if !strings.Contains(url, "www.chiphell.com") {
		url = cralwer.BASE_URL + url
	}
	doc, err := cBrowser.LoadHtmlDoc(url)
	if err != nil {
		return nil, err
	}
	return clean(doc), nil
}

func clean(doc *goquery.Document) (*Result) {
	threads := make([]Thread, 0)
	parent := doc.Find(".name").Text()
	var nextPageUrl string
	var previousPage string
	doc.Find(".page").Find("a").Each(func(i int, selection *goquery.Selection) {
		text := selection.Text()

		switch text {
		case "上一页":
			previousPage, _ = selection.Attr("href")
			if strings.Contains(previousPage, "javascript") {
				previousPage = ""
			}
		case "下一页":
			nextPageUrl, _ = selection.Attr("href")
			if strings.Contains(nextPageUrl, "javascript") {
				nextPageUrl = ""
			}
		}
	})
	doc.Find(".threadlist").Find("li").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find("a").Text()
		userName := selection.Find(".by").Text()
		num := selection.Find(".num").Text()
		url, _ := selection.Find("a").Attr("href")
		threads = append(threads, Thread{
			Title:    utils.TrimDoc(title),
			UserName: utils.TrimDoc(userName),
			Url:      utils.TrimDoc(url),
			Num:      utils.FindUInt(num),
		})
	})

	return &Result{
		Threads:      threads,
		Parent:       parent,
		PreviousPage: previousPage,
		NextPage:     nextPageUrl,
	}
}
