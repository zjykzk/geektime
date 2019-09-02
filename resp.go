package geektime

import (
	"io/ioutil"
	"net/http"
)

func readDataAndCloseResp(resp *http.Response) ([]byte, error) {
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return data, err
}
