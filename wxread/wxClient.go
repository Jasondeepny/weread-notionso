package wxread

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	WereadUrl             = "https://weread.qq.com/"
	WereadNotebooksUrl    = "https://i.weread.qq.com/user/notebooks"
	WereadBookmarklistUrl = "https://i.weread.qq.com/book/bookmarklist"
	WereadChapterInfo     = "https://i.weread.qq.com/book/chapterInfos"
	WereadReadInfoUrl     = "https://i.weread.qq.com/book/readinfo"
	WereadReviewListUrl   = "https://i.weread.qq.com/review/list"
	WereadBookInfo        = "https://i.weread.qq.com/book/info"
)

var WxCookie string

var WxClient = &http.Client{}

func LoadHomePage() {
	cookieJar, err := parseCookieString(WereadUrl, WxCookie)
	if err != nil {
		fmt.Printf("parse cookie error:%s", err)
	}
	WxClient.Jar = cookieJar
	res, _ := client(WereadUrl, "GET", nil)
	u, _ := url.Parse(WereadUrl)
	cookieJar.SetCookies(u, res.Cookies())
	WxClient.Jar = cookieJar
}

func DoWxQuery(Url string, method string, param []byte) []byte {
	_, body := client(Url, method, param)
	return body
}

func client(Url string, method string, param []byte) (*http.Response, []byte) {
	request, _ := http.NewRequest(method, Url, bytes.NewBuffer(param))
	r, err := WxClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("do request error:%s", err)
		}
	}(r.Body)
	body, _ := ioutil.ReadAll(r.Body) //
	log.Print("<<<Response StatusCode>>> :", r.StatusCode, " ,from ==>> :", Url)
	if r.StatusCode != 200 {
		log.Print("<<<WX request error>>> :", string(body))
	}
	return r, body
}

func parseCookieString(Url string, cookieString string) (*cookiejar.Jar, error) {
	var cookies []*http.Cookie
	cookieJar, _ := cookiejar.New(nil)
	cookieStrs := strings.Split(cookieString, ";")
	for _, cookie := range cookieStrs {
		parts := strings.SplitN(strings.TrimSpace(cookie), "=", 2)
		if len(parts) == 2 {
			cookie := &http.Cookie{
				Name:  parts[0],
				Value: parts[1],
			}
			cookies = append(cookies, cookie)
		}
	}
	u, _ := url.Parse(Url)
	cookieJar.SetCookies(u, cookies)
	return cookieJar, nil
}
