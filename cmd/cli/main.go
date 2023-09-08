package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/alekphilly/ReportsDownloader-go/pkg/fetch"
	"github.com/alekphilly/ReportsDownloader-go/pkg/models"
	"github.com/alekphilly/ReportsDownloader-go/pkg/parser"
)

func main() {
	ticker := flag.String("t", "", "Ticker for stock")
	flag.Parse()

	spec, _ := fetch.FetchSpec(*ticker)

	info := fetch.FetchInfo(spec["emitent_inn"].(string))

	docs := fetch.FetchDocListPage(info.Id, models.IFRSReport)

	reports := parser.ParseReports(bytes.NewReader(docs))

	fmt.Println(info.Name)
	fmt.Println(reports[0].ReportType, reports[0].ReportPeriod)

	comprRep := fetch.DownloadReport(&reports[0], "repDir")

	fmt.Println(comprRep)

	// os.Chdir("repDir")
	// r, _ := zip.OpenReader(comprRep)
	// defer r.Close()

	// for _, f := range r.File {
	// 	dst, _ := os.Create(f.Name)
	// 	compr, _ := f.Open()
	// 	io.Copy(dst, compr)
	// 	compr.Close()
	// 	dst.Close()
	// }

	os.Remove(comprRep)

	// o, _ := os.Create("out.json")
	// defer o.Close()
	// for _, t := range tr {
	// 	s, _ := json.MarshalIndent(t, "", "  ")
	// 	fmt.Fprintf(o, "%s\n", s)
	// }

}
