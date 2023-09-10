package fetch

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

const userAgent = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0"
const CompanySearchURL = "https://e-disclosure.ru/poisk-po-kompaniyam"
const companyInfoURL = "https://e-disclosure.ru/api/search/companies"
const companyDocURL = "https://e-disclosure.ru/portal/files.aspx?id=3043&type=3"

type Session struct {
	Client *http.Client
	Cookie []*http.Cookie
}

func NewSession(domain string) *Session {
	reqURL, err := url.Parse(domain)
	if err != nil {
		log.Fatal("Problem with URL: ")
	}

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", userAgent)

	client := newClient()

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return &Session{
		Client: client,
		Cookie: setCookieHeader(res.Header.Get("Set-Cookie")),
	}
}

func newClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 30 * time.Second,
	}
	return &client
}

func setCookieHeader(cookie string) []*http.Cookie {
	header := http.Header{}
	header.Add("Set-Cookie", cookie)
	req := http.Response{Header: header}
	return req.Cookies()
}
