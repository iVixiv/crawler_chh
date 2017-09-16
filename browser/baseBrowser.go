package browser

import (
	"github.com/op/go-logging"
	"github.com/sclevine/agouti"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

var (
	log = logging.MustGetLogger("browser")
)

const (
	WAITING_TIMEOUT = 10000
	WAITING_UNIT    = 300
)

// 转换 cookies 格式，将数组指针 []*http.Cookie 转换成 key=value 分号分割的字符串 cookies
func parseCookies(cookies []*http.Cookie) string {
	result := ""
	for _, cookie := range cookies {
		result += cookie.Name + "=" + cookie.Value + ";"
	}
	return result
}

// 获得 WebDriver 指针，目前使用 PhantomJS driver
func getWebDriver() *agouti.WebDriver {
	return agouti.PhantomJS()
}

func getPhoneWebDriver() *agouti.WebDriver {
	capabilities := agouti.NewCapabilities()
	capabilities["phantomjs.page.settings.userAgent"] = "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1"
	capabilitiesOption := agouti.Desired(capabilities)
	return agouti.PhantomJS(capabilitiesOption)
}

// 检查是否异常，如果异常就打印堆栈消息并结束此 goroutine
func checkErr(err error) {
	if err != nil {
		debug.PrintStack()
		log.Panic(err.Error())
	}
}

// 等待跳转到指定页面
// 死循环等待，超时时间 5s
func waiting(page *agouti.Page, redirectURL string) (mUrl string) {
	sleepTime := 0
	var url string

	for true {
		url, _ = page.URL()
		if strings.Contains(url, redirectURL) || sleepTime >= WAITING_TIMEOUT {
			break
		} else {
			sleepTime += WAITING_UNIT
			time.Sleep(time.Millisecond * WAITING_UNIT)
		}
	}
	return url
}

func wait(waitTime time.Duration) {
	time.Sleep(waitTime)
}
