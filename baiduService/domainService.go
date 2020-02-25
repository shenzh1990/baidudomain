package baiduService

import (
	"baiduDomain/utils"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type DomainHandler struct {
	AppKey    string //a48ac2b6829b481ca07793d3a58fb70c
	AppSecret string //ac6085d6d7974156b1c3f9b8d4382e55
	Url       string //https://bcd.baidubce.com
	Version   string //v1
}

const (
	Host             = "Host"
	ContentType      = "Content-Type"
	XBceDate         = "x-bce-date"
	BCE_AUTH_VERSION = "bce-auth-"
	AUTHORIZATION    = "Authorization"
	//超时时间
	DEFAULT_EXPIRATION_IN_SECONDS = "1800"
)

//验签
func (handler *DomainHandler) bceSign(req *http.Request) {
	//这样设置无效
	req.Header.Set(Host, "bcd.baidubce.com")
	req.Header.Set(ContentType, "application/json;charset=utf-8")

	timespan := time.Now().UTC().Format(time.RFC3339)
	authString := BCE_AUTH_VERSION + handler.Version + "/" + handler.AppKey +
		"/" + timespan + "/" + DEFAULT_EXPIRATION_IN_SECONDS
	req.Header.Set(XBceDate, timespan)
	//签名key
	signkey := utils.GetHmacCode(authString, handler.AppSecret)

	canonicalURI := getCanonicalURIPath(req.URL.Path)
	canonicalQueryString := getCanonicalQueryString(req.URL.Query())
	canonicalHeader, signHeader := getCanonicalHeaders(req.Header)
	CanonicalRequest := req.Method + "\n" + canonicalURI +
		"\n" + canonicalQueryString + "\n" + canonicalHeader
	signatureRequest := utils.GetHmacCode(CanonicalRequest, signkey)

	authorizationHeader := authString + "/" + signHeader + "/" + signatureRequest
	req.Header.Set(AUTHORIZATION, authorizationHeader)

}

//百度要求的uri
func getCanonicalURIPath(reqUri string) (canonicalURI string) {
	if reqUri == "" {
		return "/"
	} else if strings.HasPrefix(reqUri, "/") {
		return strings.ReplaceAll(url.QueryEscape(reqUri), "%2F", "/")
	} else {
		return "/" + strings.ReplaceAll(url.QueryEscape(reqUri), "%2F", "/")
	}
}

//百度要求的规则
func getCanonicalQueryString(query map[string][]string) string {
	keys := make([]string, 0, len(query))
	newquery := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := query[k]

		if vs[0] == AUTHORIZATION {
			continue
		}
		if vs[0] == "" {
			newquery = append(newquery, url.QueryEscape(k)+"=")
		} else {
			newquery = append(newquery, url.QueryEscape(k)+"="+url.QueryEscape(vs[0]))
		}
	}
	sort.Strings(newquery)
	return strings.Join(newquery, "&")
}
func getCanonicalHeaders(queryHeader map[string][]string) (canonicalHeaders string, signHeader string) {
	if queryHeader == nil {
		return "", ""
	}
	queryHeaderkeys := make([]string, 0, len(queryHeader))
	headerStrings := make([]string, 0, len(queryHeader))
	for k := range queryHeader {
		queryHeaderkeys = append(queryHeaderkeys, k)
	}
	sort.Strings(queryHeaderkeys)
	for _, h := range queryHeaderkeys {
		vh := queryHeader[h]
		if vh[0] == AUTHORIZATION {
			continue
		}
		headerStrings = append(headerStrings, strings.ToLower(h)+
			":"+url.QueryEscape(strings.Trim(vh[0], " ")))
	}
	sort.Strings(headerStrings)

	return strings.Join(headerStrings, "\n"),
		strings.ToLower(strings.Join(queryHeaderkeys, ";"))
}

/*
本接口用于查询域名是否可注册
*/
func (handler *DomainHandler) SearchDomain(domain string) (reslut string) {
	//response, err := http.Get(handler.Url+"/"+handler.Version+"//domain/search?domain="+domain)
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, handler.Url+"/"+handler.Version+"/domain/search?domain="+domain, nil)

	handler.bceSign(req)

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(data)
}

/*
本接口用户查询是解析记录
*/
func (handler *DomainHandler) DomianResolveList(domain []byte) (reslut []byte, err error) {
	//response, err := http.Get(handler.Url+"/"+handler.Version+"//domain/search?domain="+domain)
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, handler.Url+"/"+handler.Version+"/domain/resolve/list", bytes.NewReader(domain))

	handler.bceSign(req)

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return data, nil
}

/*
本接口用于更新解析记录
*/
func (handler *DomainHandler) DomianResolveEdit(domain []byte) (reslut bool, err error) {
	//response, err := http.Get(handler.Url+"/"+handler.Version+"//domain/search?domain="+domain)
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, handler.Url+"/"+handler.Version+"/domain/resolve/edit", bytes.NewReader(domain))

	handler.bceSign(req)
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	//data, err := ioutil.ReadAll(response.Body)
	//if err != nil {
	//	log.Fatal(err)
	//	return false,err
	//}
	//buf := new(bytes.Buffer)
	//rd,_ :=req.GetBody()
	//buf.ReadFrom(rd)
	//s := buf.String()
	//fmt.Println(
	//	handler.Url+"/"+handler.Version+"/domain/resolve/edit")
	//fmt.Println("request:POST "+req.URL.Path+" "+req.Proto)
	//if len(req.Header) > 0 {
	//	for k,v := range req.Header {
	//		fmt.Printf("%s: %s\n", k, v[0])
	//	}
	//}
	//fmt.Println("Content-Length:"+strconv.FormatInt(req.ContentLength, 10))
	//fmt.Println()
	//fmt.Println("request body:" +s)
	//fmt.Println()
	//fmt.Println(string(data))
	//fmt.Println()
	//fmt.Println(response)

	return response.StatusCode == http.StatusOK, nil
}

/*
本接口用于新增解析记录
*/
func (handler *DomainHandler) DomianResolveAdd(domain []byte) (reslut bool, err error) {
	//response, err := http.Get(handler.Url+"/"+handler.Version+"//domain/search?domain="+domain)
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, handler.Url+"/"+handler.Version+"/domain/resolve/add", bytes.NewReader(domain))

	handler.bceSign(req)

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return response.Status == string(http.StatusOK), nil
}
