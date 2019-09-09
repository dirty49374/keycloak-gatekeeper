package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type OpaRequest struct {
	Input OpaInput `json:"input"`
}

type OpaInput struct {
	Method string   `json:"method"`
	Path   []string `json:"path"`
	Token  string   `json:"token"`
}

type OpaResult struct {
	Result bool `json:"result"`
}

func checkOpaAllowed(opaURI string, resource *Resource, req *http.Request, identity *userContext) (bool, error) {

	opaRequest := OpaRequest{
		Input: OpaInput{
			Method: req.Method,
			Path:   strings.Split(req.URL.Path, "/")[1:],
			Token:  identity.token.Encode(),
		},
	}

	opaRequestBody, err := json.Marshal(opaRequest)
	if err != nil {
		return false, err
	}

	opaResponse, err := http.Post(opaURI, "application/json", bytes.NewBuffer(opaRequestBody))
	if err != nil {
		return false, err
	}
	defer opaResponse.Body.Close()

	opaResponseBody, err := ioutil.ReadAll(opaResponse.Body)
	if err != nil {
		return false, err
	}

	opaResult := OpaResult{}

	err = json.Unmarshal(opaResponseBody, &opaResult)
	if err != nil {
		return false, err
	}

	return opaResult.Result, nil
}
