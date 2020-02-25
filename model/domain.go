package model

type ResolveListIn struct {
	Domain   string `json:"domain"`
	PageNo   int    `json:"pageNo"`
	PageSize int    `json:"pageSize"`
}
type ResolveListOut struct {
	OrderBy    string `json:"orderBy"`
	Order      string `json:"order"`
	PageNo     int    `json:"pageNo"`
	PageSize   int    `json:"pageSize"`
	TotalCount int    `json:"totalCount"`
	Result     []struct {
		RecordID int    `json:"recordId"`
		Domain   string `json:"domain"`
		View     string `json:"view"`
		Rdtype   string `json:"rdtype"`
		TTL      int    `json:"ttl"`
		Rdata    string `json:"rdata"`
		ZoneName string `json:"zoneName"`
		Status   string `json:"status"`
	} `json:"result"`
}
type ResolveEditIn struct {
	Domain   string `json:"domain"`
	RdType   string `json:"rdType"`
	View     string `json:"view"`
	Rdata    string `json:"rdata"`
	TTL      int    `json:"ttl"`
	ZoneName string `json:"zoneName"`
	RecordID int    `json:"recordId"`
}

type ResolveAddIn struct {
	Domain   string `json:"domain"`
	View     string `json:"view"`
	RdType   string `json:"rdType"`
	TTL      int    `json:"ttl"`
	Rdata    string `json:"rdata"`
	ZoneName string `json:"zoneName"`
}
