package elclient

import (
	"bytes"
	"fmt"
	json "github.com/gorilla/rpc/v2/json2"
	"net/http"
)

// https://github.com/haisum/rpcexample/blob/master/examples/jrpcclient.go

// Request prepares JSON-RPC v2 request
func Request(url, method string, args interface{}) (req *http.Request, err error) {

	message, err := json.EncodeClientRequest(method, args)
	if err != nil {
		return
	}
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	return
}

// Call makes request & decodes response
func Call(req *http.Request, result interface{}) (*http.Response, error) {

	cl := new(http.Client)
	resp, err := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error in sending request to %s. %s", req.URL, err)
	}
	defer resp.Body.Close()

	err = json.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		return nil, fmt.Errorf("Couldn't decode response. %s", err)
	}
	return resp, nil
}
