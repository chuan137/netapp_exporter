package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/alecthomas/template"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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
	Name               string  `xml:"volume-id-attributes>name"`
	Owner              string  `xml:"volume-id-attributes>owning-vserver-name"`
	SizeTotal          float64 `xml:"volume-space-attributes>size-total"`
	SizeAvailable      float64 `xml:"volume-space-attributes>size-available"`
	SizeUsed           float64 `xml:"volume-space-attributes>size-used"`
	PercentageSizeUsed float64 `xml:"volume-space-attributes>percentage-size-used"`
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
const username = "admin"
const password = "netapp123"

// const url = "https://10.46.100.160/servlets/netapp.servlets.admin.XMLrequest_filer"
// const username = "mooapi"
// const password = "Api4Testing!!"

var (
	netappCapacity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netapp",
			Subsystem: "capacity",
			Name:      "svm",
			Help:      "netapp SVM capacity",
		},
		[]string{
			"svm",
			"volume",
			"metric",
		},
	)
)

func main() {
	prometheus.MustRegister(netappCapacity)

	go func() {
		for {
			getData()
			time.Sleep(60 * time.Second)
		}
	}()

	log.Print("ok")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9108", nil)
}

func getData() {
	log.Print("getData()")
	filers := queryFilers()
	for _, f := range filers {
		vols := queryVolumeByFiler(f.Name)
		// log.Print(f)
		// log.Print(vols)
		for _, v := range vols {
			netappCapacity.WithLabelValues(f.Name, v.Name, "total").Set(v.SizeTotal)
			netappCapacity.WithLabelValues(f.Name, v.Name, "available").Set(v.SizeAvailable)
			netappCapacity.WithLabelValues(f.Name, v.Name, "used").Set(v.SizeUsed)
			netappCapacity.WithLabelValues(f.Name, v.Name, "percentage_used").Set(v.PercentageSizeUsed)
		}
	}
}

func queryVolumeByFiler(filer string) []Volume {
	// post request
	xmldata, err := fetchXML(url, "./query-volume.xml", &ReqParams{250, filer})
	if err != nil {
		log.Fatal("POST query with ./query-volume.xml:", err)
	}
	v := ListVolume{}
	err = xml.Unmarshal(xmldata, &v)
	if err != nil {
		log.Print("queryVolumeByFiler(): Fail to parse xml data")
		log.Fatal(err)
	}
	return v.Volumes
}

func queryFilers() []VServer {
	// post request
	xmldata, err := fetchXML(url, "./query-vserver.xml", &ReqParams{250, ""})
	if err != nil {
		log.Fatal("POST query with ./query-vserver.xml:", err)
	}
	// decode xmldata
	v := ListVServer{}
	err = xml.Unmarshal(xmldata, &v)
	if err != nil {
		log.Print("queryFilers(): Fail to parse xml data")
		log.Fatal(err)
	}
	return v.VServers
}

func buildFromTemplate(templateFile string, params *ReqParams) *bytes.Buffer {
	var res bytes.Buffer
	tplText, err := ioutil.ReadFile(templateFile)
	if err != nil {
		log.Fatal("buildFromTemplate(): ", err)
	}
	tpl, err := template.New("vserver").Parse(string(tplText))
	if err != nil {
		log.Fatal("buildFromTemplate(): ", err)
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
