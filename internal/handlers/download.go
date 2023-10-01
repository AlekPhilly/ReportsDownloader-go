package handlers

import (
	"net/http"
	"strconv"

	"github.com/AlekPhilly/ReportsDownloader-go/internal/fetch"
)

func Download(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()["id"]
	if len(q) != 1 {
		http.Error(w, "Wrong query!!!", http.StatusBadRequest)
	}
	idx, err := strconv.Atoi(q[0])
	if err != nil || idx > len(data.Reports)-1 || idx < 0 {
		http.Error(w, "Wrong query!!!", http.StatusBadRequest)
	}
	fetch.DownloadReportToBrowser(sess, &data.Reports[idx], w)
}
