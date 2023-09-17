package parser

import (
	"io"
	"log"
	"strings"
	"time"

	"github.com/alekphilly/ReportsDownloader-go/internal/models"
	"golang.org/x/net/html"
)

func ParseReports(reportsHTML io.Reader, reportType models.DocType) []models.Report {
	rawReports := parseReportsPage(reportsHTML)
	var reports []models.Report
	var i int
	if reportType == models.AnnualReport {
		i++
	}
	for _, rep := range rawReports {
		var r models.Report
		r.ReportType = rep[1]
		r.ReportPeriod = rep[2]
		r.OriginDate, _ = time.Parse("02.01.2006", rep[i+3])
		r.PublicationDate, _ = time.Parse("02.01.2006", rep[i+4])
		r.FileLink = rep[i+5]
		reports = append(reports, r)
	}
	return reports
}

func parseReportsPage(reportsHTML io.Reader) [][]string {

	doc := html.NewTokenizer(reportsHTML)

	var isTd bool
	var rows [][]string
	var row []string

	for {
		tt := doc.Next()
		switch tt {
		case html.ErrorToken:
			err := doc.Err()
			if err == io.EOF {
				return rows
			}
			log.Fatal(err)
		case html.StartTagToken:
			t := doc.Token()
			isTd = t.Data == "td"
			if t.Data == "a" && hasAttr(t, "file-link") {
				for _, a := range t.Attr {
					if a.Key == "href" {
						row = append(row, strings.TrimSpace(a.Val))
					}
				}
			}
		case html.EndTagToken:
			t := doc.Token()
			if t.Data == "td" {
				isTd = false
			}
			if t.Data == "tr" {
				if row == nil {
					continue
				}
				if len(row) >= 6 {
					rows = append(rows, row)
				}
				row = nil
			}
			if t.Data == "table" {
				return rows
			}
		case html.TextToken:
			t := doc.Token()
			if isTd {
				if strings.TrimSpace(t.Data) != "" {
					row = append(row, strings.TrimSpace(t.Data))
				}
			}
		}
	}
}

// Check if html.Token object has an attribute
// with provided key or value
func hasAttr(n html.Token, k string) bool {
	for _, attr := range n.Attr {
		if attr.Val == k {
			return true
		} else if attr.Key == k {
			return true
		}
	}
	return false
}
