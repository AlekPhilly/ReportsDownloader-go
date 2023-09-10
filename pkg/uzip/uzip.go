package uzip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

func UnzipReport(filepath string, removeArchive bool) {
	r, _ := zip.OpenReader(filepath)
	defer r.Close()

	dir, file := path.Split(filepath)
	os.Chdir(dir)

	for i, f := range r.File {
		name := strconv.Itoa(i+1) + "-" + strings.ReplaceAll(file, ".zip", "")
		dst, err := os.Create(name)
		if err != nil {
			fmt.Println(err)
		}
		compr, _ := f.Open()
		io.Copy(dst, compr)
		compr.Close()
		dst.Close()
	}

	if removeArchive {
		os.Remove(file)
	}
}
