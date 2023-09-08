package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/alekphilly/ReportsDownloader-go/pkg/models"
)

func DownloadReport(rep *models.Report, dstDir string) string {

	cookies := getCookie(companySearchURL)

	req, err := http.NewRequest(http.MethodGet, rep.FileLink, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", userAgent)

	for _, c := range cookies {
		req.AddCookie(c)
	}

	client := newClient()

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	v := res.Header["Content-Disposition"][0]
	filename := strings.TrimSpace(strings.Replace(strings.Split(v, ";")[1], "filename=", "", 1))

	defer res.Body.Close()

	_ = os.Mkdir(dstDir, os.ModePerm)

	f, _ := os.Create(path.Join(dstDir, filename))
	defer f.Close()

	io.Copy(f, res.Body)

	return filename
}

func FetchDocListPage(companyId int, t models.DocType) []byte {

	cookies := getCookie(companySearchURL)

	docUrl, _ := url.Parse(companyDocURL)

	q := docUrl.Query()
	q.Set("id", fmt.Sprint(companyId))
	q.Set("type", fmt.Sprint(t))

	docUrl.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, docUrl.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", userAgent)

	for _, c := range cookies {
		req.AddCookie(c)
	}

	client := newClient()

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	docList := new(bytes.Buffer)
	io.Copy(docList, res.Body)

	return docList.Bytes()
}

func FetchInfo(inn string) models.CompanyInfo {

	cookies := getCookie(companySearchURL)

	form := url.Values{}
	form.Set("textfield", inn)
	form.Set("radReg", "FederalDistricts")
	form.Set("districtsCheckboxGroup", "-1")
	form.Set("regionsCheckboxGroup", "-1")
	form.Set("branchesCheckboxGroup", "-1")
	form.Set("lastPageSize", "10")
	form.Set("lastPageNumber", "1")
	form.Set("query", inn)

	req, err := http.NewRequest(http.MethodPost, companyInfoURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Accept", "application/json")

	for _, c := range cookies {
		req.AddCookie(c)
	}

	client := newClient()

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	companiesBytes := new(bytes.Buffer)
	io.Copy(companiesBytes, res.Body)

	var companiesList models.CompaniesInfoList
	err = json.Unmarshal(companiesBytes.Bytes(), &companiesList)
	if err != nil {
		log.Fatal("error unmarshalling json", err)
	}
	if len(companiesList.CompaniesList) != 1 {
		log.Fatal("Found more than one company")
	}

	return companiesList.CompaniesList[0]
}

func getCookie(domain string) []*http.Cookie {

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

	return setCookieHeader(res.Header.Get("Set-Cookie"))
}

func setCookieHeader(cookie string) []*http.Cookie {
	header := http.Header{}
	header.Add("Set-Cookie", cookie)
	req := http.Response{Header: header}
	return req.Cookies()
}
