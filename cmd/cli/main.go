package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/alekphilly/ReportsDownloader-go/pkg/fetch"
	"github.com/alekphilly/ReportsDownloader-go/pkg/models"
	"github.com/alekphilly/ReportsDownloader-go/pkg/parser"
	"github.com/alekphilly/ReportsDownloader-go/pkg/uzip"
)

func main() {
	ticker := flag.String("t", "", "Ticker for stock")
	flag.Parse()

	if len(*ticker) != 4 {
		log.Fatal("Wrong ticker!!!")
	}
	tick := strings.ToUpper(*ticker)

	spec, _ := fetch.FetchSpec(tick)

	sess := fetch.NewSession(fetch.CompanySearchURL)

	info := fetch.FetchInfo(sess, spec["emitent_inn"].(string))

	docs := fetch.FetchDocListPage(sess, info.Id, models.IFRSReport)

	reports := parser.ParseReports(bytes.NewReader(docs))

	fmt.Println(info.Name)
	fmt.Println(reports[0].ReportType, reports[0].ReportPeriod)

	comprRep := fetch.DownloadReport(sess, &reports[0], "repDir")

	fmt.Println(comprRep)

	uzip.UnzipReport(comprRep, true)

	// o, _ := os.Create("out.json")
	// defer o.Close()
	// for _, t := range tr {
	// 	s, _ := json.MarshalIndent(t, "", "  ")
	// 	fmt.Fprintf(o, "%s\n", s)
	// }

}
