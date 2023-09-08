package fetch

import (
	"net/http"
	"time"
)

const userAgent = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0"
const companySearchURL = "https://e-disclosure.ru/poisk-po-kompaniyam"
const companyInfoURL = "https://e-disclosure.ru/api/search/companies"
const companyDocURL = "https://e-disclosure.ru/portal/files.aspx?id=3043&type=3"

func newClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 30 * time.Second,
	}
	return &client
}
