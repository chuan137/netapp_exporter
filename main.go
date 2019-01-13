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
	Filer      string
}

// VServer is type for VServer
type VServer struct {
	Name string `xml:"vserver-name"`
	UUID string `xml:"uuid"`
	// IPSpace string `xml:"ipspace"`
}

// Volume is type for Volume
type Volume struct {
	Name          string `xml:"volume-id-attributes>name"`
	Owner         string `xml:"volume-id-attributes>owning-vserver-name"`
	SizeAvailable string `xml:"volume-space-attributes>size-available"`
	SizeUsed      string `xml:"volume-space-attributes>size-used"`
}

// ListVServer is type for list of VServers
type ListVServer struct {
	XMLName       xml.Name  `xml:"netapp"`
	NetappVersion string    `xml:"version,attr"`
	NumRec        int       `xml:"results>num-records"`
	VServers      []VServer `xml:"results>attributes-list>vserver-info"`
}

// ListVolume is type for list of Volumes
type ListVolume struct {
	XMLName       xml.Name `xml:"netapp"`
	NetappVersion string   `xml:"version,attr"`
	NumRec        int      `xml:"results>num-records"`
	Volumes       []Volume `xml:"results>attributes-list>volume-attributes"`
}

const url = "https://10.44.58.21/servlets/netapp.servlets.admin.XMLrequest_filer"

func main() {
	filers := queryFilers()

	for _, f := range filers {
		log.Print(f)
		vols := queryVolumeByFiler(f.Name)
		log.Print(vols)
	}

	// http.Handle("/metrics", promhttp.Handler())
	// http.ListenAndServe("localhost:8080", nil)
}

func queryVolumeByFiler(filer string) []Volume {
	// post request
	xmldata, err := fetchXML(url, "./query-volume.xml", &ReqParams{250, filer})
	if err != nil {
		log.Fatal(err)
	}
	v := ListVolume{}
	err = xml.Unmarshal(xmldata, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v.Volumes
}

func queryFilers() []VServer {
	// post request
	xmldata, err := fetchXML(url, "./query-vserver.xml", &ReqParams{250, ""})
	if err != nil {
		log.Fatal(err)
	}
	// decode xmldata
	v := ListVServer{}
	err = xml.Unmarshal(xmldata, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v.VServers
}

func buildFromTemplate(templateFile string, params *ReqParams) *bytes.Buffer {
	var res bytes.Buffer
	tplText, _ := ioutil.ReadFile(templateFile)
	tpl, err := template.New("vserver").Parse(string(tplText))
	if err != nil {
		log.Fatal(err)
	}
	tpl.Execute(&res, params)
	return &res
}

func fetchXML(url string, reqTemplateFile string, reqParams *ReqParams) ([]byte, error) {
	// fill parameters into template
	xmlbody := buildFromTemplate(reqTemplateFile, reqParams)
	// request payload
	req, err := http.NewRequest("POST", url, xmlbody)
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
