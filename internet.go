package internet

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type Internet struct {
	lastUrl  string
	client   *http.Client
	header   *http.Header
	response *http.Response
}

func New() Internet {
	jar, _ := cookiejar.New(nil)
	e := Internet{
		"",
		&http.Client{
			Jar: jar,
			//			CheckRedirect: func(req *http.Request, via []*http.Request) error {
			//				return http.ErrUseLastResponse
			//			},
		},
		&http.Header{},
		nil,
	}
	return e
}

func (i Internet) AddCookie(scheme string, host string, key string, value string) {
	i.client.Jar.SetCookies(&url.URL{Scheme: scheme, Host: host}, []*http.Cookie{&http.Cookie{Name: key, Value: value}})
}
func (i Internet) FlushCookies() {
	i.client.Jar.Cookies(&url.URL{Scheme: "http", Host: "www.google.com"})
}
func (i Internet) AddHeader(key string, value string) {
	i.header.Add(key, value)
}
func (i Internet) RemoveHeader(key string) {
	i.header.Del(key)
}
func (i Internet) PrintRequestHeaders() {
	log.Println(i.header)
}
func (i Internet) PrintResponseHeaders() {
	if i.response != nil {
		log.Println(i.response.Header)
	}
}
func (i Internet) PrintCookies(scheme string, host string) {
	for _, cookie := range i.client.Jar.Cookies(&url.URL{Scheme: scheme, Host: host}) {
		log.Println(cookie)
	}
}

func (i Internet) GetUrl(targetUrl string) ([]byte, error) {
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Safari/537.36")
	req.Header.Set("Referer", i.lastUrl)
	i.response, err = i.client.Do(req)
	i.lastUrl = targetUrl
	if err != nil {
		log.Fatalln(err)
	}
	defer i.response.Body.Close()
	body, err := ioutil.ReadAll(i.response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body, err
}

func (i Internet) PostUrl(targetUrl string, payload string) ([]byte, error) {

	data := url.Values{}
	for _, v := range strings.Split(payload, "&") {
		data.Set(strings.Split(v, "=")[0], strings.Split(v, "=")[1])
	}

	req, err := http.NewRequest("POST", targetUrl, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Safari/537.36")
	req.Header.Set("Referer", i.lastUrl)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	i.response, err = i.client.Do(req)
	i.lastUrl = targetUrl
	if err != nil {
		log.Fatalln(err)
	}

	defer i.response.Body.Close()
	body, err := ioutil.ReadAll(i.response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body, err
}
