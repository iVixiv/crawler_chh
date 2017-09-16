package browser

import "github.com/PuerkitoBio/goquery"

type Browser interface {
	LoadHtmlDoc(url string) (*goquery.Document, error)
}
