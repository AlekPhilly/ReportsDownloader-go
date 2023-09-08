package models

import "time"

type DocType int

const (
	AnnualReport DocType = iota + 2
	AccountReport
	IFRSReport
	IssuerReport
)

type Report struct {
	ReportType      string    `json:"report_type"`
	ReportPeriod    string    `json:"report_period"`
	OriginDate      time.Time `json:"origin_date"`
	PublicationDate time.Time `json:"publication_date"`
	FileLink        string    `json:"file_link"`
}
