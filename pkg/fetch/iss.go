package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type securitySpec struct {
	Securities struct {
		Columns []string `json:"columns"`
		Data    [][]any  `json:"data"`
	} `json:"securities"`
}

func cleanSpec(ticker string, spec *securitySpec) map[string]any {
	var fields []any
	for _, dataSet := range spec.Securities.Data {
		for _, field := range dataSet {
			if field == ticker {
				fields = dataSet
				break
			}
		}
		if fields != nil {
			break
		}
	}

	shareMap := make(map[string]any)
	for i := 0; i < len(fields); i++ {
		shareMap[spec.Securities.Columns[i]] = fields[i]
	}

	return shareMap
}

func FetchSpec(ticker string) (map[string]any, error) {
	req, _ := http.NewRequest(http.MethodGet, specURL(ticker), nil)

	res, err := requestSpec(req)
	if err != nil {
		log.Fatal(err)
	}

	var spec *securitySpec
	err = json.Unmarshal(res, &spec)
	if err != nil {
		fmt.Printf("error decoding json: %s\n", err)
	}
	// fmt.Printf("%v", spec)

	return cleanSpec(ticker, spec), nil
}

func requestSpec(req *http.Request) ([]byte, error) {

	client := newClient()
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body := new(bytes.Buffer)
	_, err = io.Copy(body, res.Body)
	if err != nil {
		fmt.Printf("error reading response body: %s\n", err)
		return nil, err
	}
	return body.Bytes(), nil
}

func specURL(ticker string) string {

	specReq, err := url.Parse("https://iss.moex.com/iss/securities.json")
	if err != nil {
		log.Fatal(err)
	}
	q := specReq.Query()
	q.Set("q", ticker)
	specReq.RawQuery = q.Encode()

	return specReq.String()
}
