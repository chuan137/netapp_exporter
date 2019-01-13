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

const url = "https://10.44.58.21/servlets/netapp.servlets.admin.XMLrequest_filer"

func main() {
	filers := getFilers()

	for _, f := range filers {
		log.Print(f)
	}
	// http.Handle("/metrics", promhttp.Handler())
	// http.ListenAndServe(":8080", nil)
}

func getFilers() []VServer {
	// post request
	xmldata, err := fetchXML(url, "./query-vserver.xml", &ReqParams{250, ""})
	if err != nil {
		log.Fatal(err)
	}
	// decode xmldata
	v := ResVServer{}
	err = xml.Unmarshal(xmldata, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v.VServers
}

func fetchXML(url string, reqTemplateFile string, reqParams *ReqParams) ([]byte, error) {
	// request template
	tplText, _ := ioutil.ReadFile(reqTemplateFile)
	tpl, err := template.New("vserver").Parse(string(tplText))
	if err != nil {
		return []byte{}, err
	}
	// fill parameters into template
	var xmlbody bytes.Buffer
	tpl.Execute(&xmlbody, reqParams)
	// request payload
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
	// extract response data
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}
