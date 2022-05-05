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
	lastUrl string
	client  *http.Client
}

func New() Internet {
	jar, _ := cookiejar.New(nil)
	e := internet{
		"",
		&http.Client{
			Jar: jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
	return e
}

func (i Internet) GetUrl(targetUrl string) ([]byte, error) {
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Safari/537.36")
	req.Header.Set("Referer", i.lastUrl)
	resp, err := i.client.Do(req)
	i.lastUrl = targetUrl
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
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

	resp, err := i.client.Do(req)
	i.lastUrl = targetUrl
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body, err
}
