package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alecthomas/template"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
	<netapp vfiler="abc" version="1.7" xmlns="http://www.netapp.com/filer/admin">
*/

// ReqParams is type of request parameters
type ReqParams struct {
	MaxRecords int
	VFiler     string
}

func main() {
	const url = "https://10.44.58.21/servlets/netapp.servlets.admin.XMLrequest_filer"

	var xmlGetVserver bytes.Buffer
	tplText, _ := ioutil.ReadFile("./query-vserver.xml")
	tpl := template.Must(template.New("vserver").Parse(string(tplText)))
	tpl.Execute(&xmlGetVserver, &ReqParams{250, ""})

	// result, err := fetchXML(url, xmlGetVserver.Bytes())
	result, err := fetchXML(url, string(tplText), &ReqParams{250, ""})
	if err != nil {
		log.Fatal(err)
	}
	log.Print(result)

	// http.Handle("/metrics", promhttp.Handler())
	// http.ListenAndServe(":8080", nil)
}

func fetchXML(url string, reqTemplate string, reqParams *ReqParams) (string, error) {
	var xmlbody bytes.Buffer

	// parse template
	tpl, err := template.New("vserver").Parse(reqTemplate)
	if err != nil {
		return "", err
	}
	tpl.Execute(&xmlbody, reqParams)

	// new request
	req, err := http.NewRequest("POST", url, &xmlbody)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("admin", "netapp123")
	req.Header.Add("Content-Type", "text/xml")

	// do request
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// parse result
	log.Print(resp.Status)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
