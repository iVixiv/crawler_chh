package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/golibs/uuid"
	"github.com/op/go-logging"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	log        = logging.MustGetLogger("utils")
	regFloat64 = regexp.MustCompile("[1-9]\\d*.\\d*|0.\\d*[1-9]\\d*")
	regInt     = regexp.MustCompile("[1-9]\\d*")
	regDate    = regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}")
)

const (
	CHINA_TIME_LAYOUT   = "2006年01月02日"
	DEFAULT_TIME_LAYOUT = "2006-01-02"
	PRIMARY_TIME_LAYOUT = "2006-01-02 15:04:05"
)

// 删除任何形式的空符号，包括 \t \n " "
func TrimAnythingEmpty(str string) string {
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, " ", "", -1)
	return str
}

func StringTimeToTimeByChinaLayout(timeStr string) (time.Time, error) {
	return StringTimeToTime(CHINA_TIME_LAYOUT, timeStr)
}

func StringTimeToTimeByRFC3339Layout(timeStr string) (time.Time, error) {
	return StringTimeToTime(time.RFC3339, timeStr)
}

// 字符串类型的时间转换成 time.Time 使用默认 layout ""
func StringTimeToTimeByDefaultLayout(timeStr string) (time.Time, error) {
	return StringTimeToTime(DEFAULT_TIME_LAYOUT, timeStr)
}

// 用正则表达式匹配 yyyy-MM-dd 格式的日期
func FindDefaultLayoutTime(str string) (time.Time, error) {
	timeStr := regDate.FindString(str)
	return StringTimeToTimeByDefaultLayout(timeStr)
}

func StringTimeToTimeByPrimaryTimeLayout(timeStr string) (time.Time, error) {
	return StringTimeToTime(PRIMARY_TIME_LAYOUT, timeStr)
}

// 字符串类型的时间转换成 time.Time
func StringTimeToTime(layout string, timeStr string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation(layout, timeStr, loc)
}

// 时间戳转换到 time.Time
func UnixTimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// 计算实际借款天数
func BorrowDays(borrowTime time.Time) int {
	now := time.Now()
	second := now.Unix() - borrowTime.Unix()
	oneDaySecond := int64(60 * 60 * 24)
	return int(second/oneDaySecond) + 1
}

// 提取字符串中的正浮点数
func FindFloat64(str string) float64 {
	money := regFloat64.FindString(str)
	if money == "" {
		return 0
	}
	money = strings.Replace(money, ",", "", -1) // 去掉金额中的逗号
	float, err := strconv.ParseFloat(money, 64)
	if err != nil {
		log.Warning("Find float64 [%s] error %s", str, err.Error())
	}
	return float
}

// 提取字符串中的正整数
func FindUInt(str string) uint {
	if str == "" {
		return 0
	}
	integer, err := strconv.ParseInt(regInt.FindString(str), 10, 32)
	if err != nil {
		log.Warning("Find uint [%s] error: %s", str, err.Error())
	}
	return uint(integer)
}

func FindUInt8(str string) uint8 {
	if str == "" {
		return 0
	}
	integer, err := strconv.ParseInt(regInt.FindString(str), 10, 16)
	if err != nil {
		log.Warning("Find uint8 [%s] error: %s", str, err.Error())
	}
	return uint8(integer)
}

func MD5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// 计算 ID，有些贷款平台对于账单没有一个唯一的 ID 返回，我们将用户名和申请时间/还款时间作为这条订单的唯一标识
// MD5(platform + phone + datestr)
func ComputeIdByPhoneAndDateStr(platform string, phone string, datestr string) string {
	return MD5(platform + phone + datestr)
}

// date layout 2006-01-02
// 弃用，建议用 ComputeIdByPhoneAndDatePrimaryLayout 方法
func ComputeIdByPhoneAndDate(platform string, phone string, date time.Time) string {
	return ComputeIdByPhoneAndDateStr(platform, phone, date.Format(DEFAULT_TIME_LAYOUT))
}

// date layout 2006-01-02 15:04:05
func ComputeIdByPhoneAndDatePrimaryLayout(platform string, phone string, date time.Time) string {
	return ComputeIdByPhoneAndDateStr(platform, phone, date.Format(PRIMARY_TIME_LAYOUT))
}

// 生成一个爬取唯一标识，使用 UUID 算法
func CrawlIdGen() string {
	rand := uuid.Rand()
	x := [16]byte(rand)
	return fmt.Sprintf("%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x",
		x[0], x[1], x[2], x[3], x[4],
		x[5], x[6],
		x[7], x[8],
		x[9], x[10], x[11], x[12], x[13], x[14], x[15])
}

// 单位转换，分 -> 元
func FenToYuan(money int) float64 {
	return float64(money / 100)
}

// 已还金额，单期账单的已还金额是根据还款状态
// 如果已还款，已还金额等于借款总金额，否则等于 0
func HadRepayAmount(status bool, amount float64) float64 {
	if status {
		return amount
	} else {
		return 0
	}
}

// 已还期数，单期账单的已还期数是根据还款状态
// 如果已还款，已还期数等于1，否则等于 0
func HadRepayPeriod(status bool) uint8 {
	if status {
		return 1
	} else {
		return 0
	}
}

// 单期账单总期数
func SingleTotalPeriod() uint8 {
	return 1
}

// 获取Authorization值
func SetAuthorization(userPhone string, token string, req *http.Request, authorization string) {
	str := userPhone + ":" + token
	req.Header.Add(authorization, "Basic "+base64.StdEncoding.EncodeToString([]byte(str)))
}

// 获取设置了种子的 *rand.Rand
func RandSeed() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// 获取当前时间的日期格式 2006-01-02 15:04:05 yyyy-MM-dd HH:mm:ss
func GetCurrentTimeDateType() string {
	return time.Now().Format(PRIMARY_TIME_LAYOUT)
}

//字符串截取 从doc中找出第一个 start字符开始 到end结束( 不包含start)
func SplitDoc(doc string, start string, end string) string {
	n := strings.Index(doc, start)
	if n == -1 {
		n = 0
	} else {
		n += len(start)
	}
	doc = string([]byte(doc)[n:])
	m := strings.Index(doc, end)
	if m == -1 {
		m = len(doc)
	}
	doc = string([]byte(doc)[:m])
	return doc
}

//字符串去空格和换行
func TrimDoc(doc string) string {
	// 去除空格
	doc = strings.Replace(doc, " ", "", -1)
	// 去除换行符
	doc = strings.Replace(doc, "\n", "", -1)
	return doc
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
