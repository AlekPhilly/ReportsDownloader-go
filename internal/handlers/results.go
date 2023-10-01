package handlers

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
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
	tick := strings.ToUpper(r.FormValue("ticker"))
	repType, _ := strconv.Atoi(r.FormValue("reportType"))
	rt := models.DocType(repType)

	spec, _ := fetch.FetchSpec(tick)

	inn, ok := spec["emitent_inn"]
	if !ok {
		http.Error(w, "Can't fetch company info", http.StatusNotFound)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		panic(http.ErrAbortHandler)
	}

	sess = fetch.NewSession(fetch.CompanySearchURL)

	info := fetch.FetchInfo(sess, inn.(string))

	docs := fetch.FetchDocListPage(sess, info.Id, rt)

	reports := parser.ParseReports(bytes.NewReader(docs), rt)

	data = Data{
		Info:    info,
		Reports: reports,
	}

	files := []string{
		"./html/results.page.tmpl",
		"./html/base.layout.tmpl",
	}

	resPage := template.Must(template.ParseFiles(files...))

	resPage.Execute(w, data)
}
