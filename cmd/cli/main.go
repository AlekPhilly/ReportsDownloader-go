package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/alekphilly/ReportsDownloader-go/pkg/fetch"
	"github.com/alekphilly/ReportsDownloader-go/pkg/models"
	"github.com/alekphilly/ReportsDownloader-go/pkg/parser"
	"github.com/alekphilly/ReportsDownloader-go/pkg/uzip"
)

const repDir = "reports"

func main() {
	ticker := flag.String("s", "", "Ticker for stock")
	repType := flag.Int("t", 3, "Report type: 1 - Annual, 2 - Account, 3 - IFRS, 4 - Issuer")
	flag.Parse()

	r := *repType
	if r < 1 || r > 4 {
		log.Fatal("Wrong report type!!!")
	}
	rt := models.DocType(r + 1)

	if len(*ticker) != 4 {
		log.Fatal("Wrong ticker!!!")
	}
	tick := strings.ToUpper(*ticker)

	spec, _ := fetch.FetchSpec(tick)

	sess := fetch.NewSession(fetch.CompanySearchURL)

	info := fetch.FetchInfo(sess, spec["emitent_inn"].(string))

	docs := fetch.FetchDocListPage(sess, info.Id, rt)

	reports := parser.ParseReports(bytes.NewReader(docs), rt)

	fmt.Println(info.Name)
	fmt.Println(reports[0].ReportType, reports[0].ReportPeriod)

	comprRep := fetch.DownloadReport(sess, &reports[0], path.Join(repDir, info.Name))

	fmt.Println("Saved to:", comprRep)

	uzip.UnzipReport(comprRep, true)
}
