package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alecthomas/template"
)

func fetchXML(url, username, password string, xmlbody *bytes.Buffer) ([]byte, error) {
	// request payload
	req, err := http.NewRequest("POST", url, xmlbody)
	if err != nil {
		return []byte{}, err
	}
	req.SetBasicAuth(username, password)
	req.Header.Add("Content-Type", "text/xml")

	// do request
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	// extract response data
	if resp.Status != "200 OK" {
		log.Fatal(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func makeRequestFromTemplateFile(fname string, params *ReqParams) *bytes.Buffer {
	// queryFile string, queryParams *ReqParams
	tpl, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	return makeRequestFromTemplate(tpl, params)
}

func makeRequestFromTemplate(tpl []byte, params interface{}) *bytes.Buffer {
	var res bytes.Buffer
	t, err := template.New("__query").Parse(string(tpl))
	if err != nil {
		log.Fatal("buildFromTemplate(): ", err)
	}
	t.Execute(&res, params)
	return &res
}
