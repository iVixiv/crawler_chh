package home

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
	"github.com/op/go-logging"
	"crawler_chh/utils"
	"crawler_chh/cralwer"
	"net/http"
	"io/ioutil"
	"fmt"
	"crawler_chh/browser"
)

const (
	HOME = "forum.php"
)

var (
	c_browser = new(browser.Chh_browser)
	log       = logging.MustGetLogger("chh_home")
)

func (this Home) Crawler() (*[]Result, error) {
	doc, err := c_browser.LoadHtmlDoc(cralwer.BASE_URL + HOME)
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

func fetchHtmlDOC(url string) (*goquery.Document, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("accept-encoding", "gzip, deflate, br")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.8,en;q=0.6")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("referer", "gzip, deflate, br")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1")
	req.Header.Add("cookie", "v2x4_48dd_saltkey=Dz9oezhE; v2x4_48dd_lastvisit=1504320270; v2x4_48dd_visitedfid=80D53D170D297D188D36; rcount=c4b0a56a75122c9d536a3dde17f025bb; ruser=94937; PHPSESSID=5bc6bfc91c12f72748f27df8625c1f5c; v2x4_48dd_st_p=0%7C1505699131%7C928f5c1c2db403845b96a31a6bd7b777; v2x4_48dd_viewid=tid_1775150; v2x4_48dd_st_t=0%7C1505699139%7C8e7afce16f6f5a8197ab28912e12474f; v2x4_48dd_forum_lastvisit=D_36_1505449614D_188_1505449644D_297_1505449675D_170_1505545013D_53_1505549280D_80_1505699139; v2x4_48dd_sendmail=1; _ga=GA1.2.265181952.1504323873; _gat=1; v2x4_48dd_lastact=1505704307%09forum.php%09; v2x4_48dd_sid=HNpI65; _fmdata=F1EA6790D4849B0A9530211CCA333CDE25E1FBE858E208E51CE4585E413E7C2AAFB7DCC56EDFE5D38E9558D31A39C9DD772D69FDD54A5F98")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	//if mahonia.GetCharset("utf-8") == nil {
	//	return nil, fmt.Errorf("charset not suported \n")
	//}
	//dec := mahonia.NewDecoder("gbk")

	//doc := dec.ConvertString(string(body))
	fmt.Printf("%+v", res.Header)
	fmt.Println(string(body))

	return goquery.NewDocumentFromReader(strings.NewReader(string(body)))
}
