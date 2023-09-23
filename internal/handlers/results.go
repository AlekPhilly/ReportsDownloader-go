package handlers

import (
	"bytes"
	"net/http"
	"strconv"
	"text/template"

	"github.com/AlekPhilly/ReportsDownloader-go/internal/fetch"
	"github.com/AlekPhilly/ReportsDownloader-go/internal/models"
	"github.com/AlekPhilly/ReportsDownloader-go/internal/parser"
)

type Data struct {
	Info    models.CompanyInfo
	Reports []models.Report
}

var data Data
var sess *fetch.Session

func Search(w http.ResponseWriter, r *http.Request) {
	tick := r.FormValue("ticker")
	repType, _ := strconv.Atoi(r.FormValue("reportType"))
	rt := models.DocType(repType)

	spec, _ := fetch.FetchSpec(tick)

	sess = fetch.NewSession(fetch.CompanySearchURL)

	info := fetch.FetchInfo(sess, spec["emitent_inn"].(string))

	docs := fetch.FetchDocListPage(sess, info.Id, rt)

	reports := parser.ParseReports(bytes.NewReader(docs), rt)

	data = Data{
		Info:    info,
		Reports: reports,
	}

	resPage := template.Must(template.ParseFiles("./html/tmpl/results.html"))

	resPage.Execute(w, data)
}

func Download(w http.ResponseWriter, r *http.Request) {
	idx, _ := strconv.Atoi(r.FormValue("radio"))
	fetch.DownloadReportToBrowser(sess, &data.Reports[idx], w)
}
