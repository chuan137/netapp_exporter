package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
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

// VServer is type for VServer
type VServer struct {
	Name string `xml:"vserver-name"`
	UUID string `xml:"uuid"`
	// IPSpace string `xml:"ipspace"`
}

// ResVServer is type for VServer of all results
type ResVServer struct {
	XMLName       xml.Name  `xml:"netapp"`
	NetappVersion string    `xml:"version,attr"`
	NumRec        int       `xml:"results>num-records"`
	VServers      []VServer `xml:"results>attributes-list>vserver-info"`
}

func main() {
	const url = "https://10.44.58.21/servlets/netapp.servlets.admin.XMLrequest_filer"

	tplText, _ := ioutil.ReadFile("./query-vserver.xml")
	result, err := fetchXML(url, string(tplText), &ReqParams{250, ""})
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(result))

	parseVServer(result)

	// http.Handle("/metrics", promhttp.Handler())
	// http.ListenAndServe(":8080", nil)
}

func parseVServer(xmldata []byte) {
	v := ResVServer{}
	err := xml.Unmarshal(xmldata, &v)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(v)

	for _, s := range v.VServers {
		log.Print(s)
	}
}

func fetchXML(url string, reqTemplate string, reqParams *ReqParams) ([]byte, error) {
	var xmlbody bytes.Buffer

	// parse template
	tpl, err := template.New("vserver").Parse(reqTemplate)
	if err != nil {
		return []byte{}, err
	}
	tpl.Execute(&xmlbody, reqParams)

	// new request
	req, err := http.NewRequest("POST", url, &xmlbody)
	if err != nil {
		return []byte{}, err
	}
	req.SetBasicAuth("admin", "netapp123")
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

	// parse result
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}
