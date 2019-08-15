package network

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

func GetLocalPublicIp() (string, error) {
	// `nc ns1.dnspod.cn 6666`
	res, err := http.Get("http://ifconfig.me/ip")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	reg := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	return reg.FindString(string(result)), nil
}
