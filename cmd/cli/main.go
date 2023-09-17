package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
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
	repNum := flag.Int("n", 5, "How many reports to scan")
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

	maxRepListLength := *repNum

	spec, _ := fetch.FetchSpec(tick)

	sess := fetch.NewSession(fetch.CompanySearchURL)

	info := fetch.FetchInfo(sess, spec["emitent_inn"].(string))

	docs := fetch.FetchDocListPage(sess, info.Id, rt)

	reports := parser.ParseReports(bytes.NewReader(docs), rt)

	fmt.Println(info.Name)
	fmt.Println()
	fmt.Println("Available reports:")

	var repCount int
	for i, rep := range reports {
		fmt.Println(i+1, "-", rep.ReportType, rep.ReportPeriod)
		repCount++
		if repCount == maxRepListLength {
			fmt.Println()
			break
		}
	}

	fmt.Println("Which reports do you want to download? (numbers separated by space)")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	dlList := scanner.Text()

	for _, n := range strings.Fields(dlList) {
		i, err := strconv.Atoi(n)
		if err != nil {
			fmt.Println("Wrong report number:", n)
			continue
		}
		i--

		comprRep := fetch.DownloadReport(sess, &reports[i], path.Join(repDir, info.Name))
		fmt.Printf("Report %d saved to %s\n", i, comprRep)
		uzip.UnzipReport(comprRep, true)
	}
}
