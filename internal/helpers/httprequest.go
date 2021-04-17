package helpers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

//HTTPRequest is a generic template for making http requests
func HTTPRequest(client *http.Client, url string, method string, payload []byte) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		err = fmt.Errorf("[ERROR] HTTPRequest: creating request \n%s", err.Error())
		return nil, 0, err
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("[ERROR] HTTPRequest: doing request \n%s", err.Error())
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("[ERROR] HTTPRequest: reading response body \n%s", err.Error())
		return nil, 0, err
	}
	return body, resp.StatusCode, nil
}
