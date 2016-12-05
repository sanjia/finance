package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func HttpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("find exception:%s\r\n", err.Error())
		return "", err
	}
	return string(body), nil
}
