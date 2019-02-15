package netapp

import (
	"encoding/xml"
)

type Base struct {
	XMLName xml.Name `xml:"netapp"`
	Version string   `xml:"version,attr,omitempty"`
	Vfiler  string   `xml:"vfiler,attr,omitempty"`
	Client
}

func (b *Base) setVersion(v string) *Base {
	b.Version = v
	return b
}

// func (b *Base) get() {
// 	b.Client.do()
// }
