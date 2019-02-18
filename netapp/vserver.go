package netapp

import (
	"encoding/xml"
	"log"
)

type Vserver struct {
	Base
	Params struct {
		XMLName           xml.Name
		MaxRecords        int         `xml:"max-records,omitempty"`
		Query             interface{} `xml:"query>vserver-info,omitempty"`
		DesiredAttribtues interface{} `xml:"desired-attributes>vserver-info,omitempty"`
	}
}

// type VserverInfo struct {
// 	VserverName string `xml:"vserver-name,omitempty"`
// 	UUID        string `xml:"uuid,omitempty"`
// }

type VserverName struct {
	V string `xml:"vserver-name"`
}

type UUID struct {
	V string `xml:"uuid"`
}

type VserverType struct {
	V string `xml:"vserver-type"`
}

func NewVserver(attr interface{}) *Vserver {
	v := &Vserver{}
	v.Base = Base{Version: "1.7"}
	// v.Params.MaxRecords = 100
	v.Params.DesiredAttribtues = attr
	return v
}

func (v *Vserver) List() {
	v.Params.XMLName.Local = "vserver-get-iter"

	result := v.Base.Client.do(&v)
	log.Print(result)
}
