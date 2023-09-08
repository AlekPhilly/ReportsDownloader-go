package models

type CompanyInfo struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	District     string `json:"district"`
	Region       string `json:"region"`
	Branch       string `json:"branch"`
	LastActivity string `json:"lastActivity"`
	DocCount     int    `json:"docCount"`
}

type CompaniesInfoList struct {
	CompaniesList []CompanyInfo `json:"foundCompaniesList"`
}
