package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"
)

func GetHmacCode(value string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

/**
* type ==4 获取ipv4
* type ==6 获取ipv6
* type ==test  v4优先返回v4地址 否则返回v6
**/
func GetIp(t string) (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, "https://"+t+".ipw.cn", nil)

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
