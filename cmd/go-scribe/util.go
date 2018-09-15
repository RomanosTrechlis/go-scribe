package main

import (
	"strings"
	"net/http"
	"io/ioutil"
)

func readConfigurationFile(file string) ([]byte, error) {
	if strings.HasPrefix(file, "https") || strings.HasPrefix(file, "http") {
		resp, err := http.Get(file)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		return ioutil.ReadAll(resp.Body)
	}

	return ioutil.ReadFile(file)
}
