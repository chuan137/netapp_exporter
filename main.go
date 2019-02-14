package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/alecthomas/template"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

// Parameter
var (
	waitTime = kingpin.Flag("wait", "Wait time").Short('w').Default("300").Int()
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

// VolumeList is type for list of Volumes
type VolumeList struct {
	XMLName       xml.Name `xml:"netapp"`
	NetappVersion string   `xml:"version,attr"`
	NumRec        int      `xml:"results>num-records"`
	Volumes       []Volume `xml:"results>attributes-list>volume-attributes"`
}

type filerConfig struct {
	Name     string `yaml:"name"`
	IP       string `yaml:"ip"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

const urlTemplate = "https://%s/servlets/netapp.servlets.admin.XMLrequest_filer"

var (
	netappCapacity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netapp",
			Subsystem: "capacity",
			Name:      "svm",
			Help:      "netapp SVM capacity",
		},
		[]string{
			"filer",
			"svm",
			"volume",
			"metric",
		},
	)
)

func main() {

	filers := readFilerConfig("./netapp_filers.yaml")

	prometheus.MustRegister(netappCapacity)

	go func() {
		for {
			for _, f := range filers {
				getData(&f)
			}
			time.Sleep(time.Duration(*waitTime) * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9108", nil)
}

func getData(fc *filerConfig) {
	log.Print("[Info] getData() ", fc.Name)

	svms := querySvmByFiler(fc)

	for _, vs := range svms {
		vols := queryVolumeByFiler(fc, &vs)
		// log.Print(s)
		// log.Print(vols)
		for _, v := range vols {
			netappCapacity.WithLabelValues(fc.Name, vs.Name, v.Name, "total").Set(v.SizeTotal)
			netappCapacity.WithLabelValues(fc.Name, vs.Name, v.Name, "available").Set(v.SizeAvailable)
			netappCapacity.WithLabelValues(fc.Name, vs.Name, v.Name, "used").Set(v.SizeUsed)
			netappCapacity.WithLabelValues(fc.Name, vs.Name, v.Name, "percentage_used").Set(v.PercentageSizeUsed)
		}
	}
}

func queryVolumeByFiler(fc *filerConfig, vs *VServer) (v []Volume) {
	url := fmt.Sprintf(urlTemplate, fc.IP)

	// post request
	xmldata, err := fetchXML(url, fc.Username, fc.Password, "./query_volume.xml", &ReqParams{250, vs.Name})
	if err != nil {
		log.Print("[Warning] ", err)
		return
	}

	var l VolumeList
	err = xml.Unmarshal(xmldata, &l)
	if err != nil {
		log.Print("queryVolumeByFiler(): Fail to parse xml data")
		log.Fatal(err)
	}

	return l.Volumes
}

func querySvmByFiler(fc *filerConfig) (v []VServer) {
	url := fmt.Sprintf(urlTemplate, fc.IP)

	// post request
	xmldata, err := fetchXML(url, fc.Username, fc.Password, "./query_vserver.xml", &ReqParams{250, ""})
	if err != nil {
		log.Print("[Warning] ", err)
		return
	}

	// decode xmldata
	var l ListVServer
	err = xml.Unmarshal(xmldata, &l)
	if err != nil {
		log.Print("queryFilers(): Fail to parse xml data")
		log.Fatal(err)
	}
	return l.VServers
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

func fetchXML(url, username, password string, reqTemplateFile string, reqParams *ReqParams) ([]byte, error) {
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

func readFilerConfig(fileName string) (c []filerConfig) {

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("[ERROR] ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatal("[ERROR] ", err)
	}

	return
}
