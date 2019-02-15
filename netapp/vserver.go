package netapp

type Vserver struct {
	Base
	DesiredAttribtues interface{} `xml:"desired-attributes>vserver-info"`
}

type VserverInfo struct {
	VserverName string `xml:"vserver-name,omitempty"`
	UUID        string `xml:"uuid,omitempty"`
}

type VserverName struct {
	Value string `xml:"vserver-name"`
}

type UUID struct {
	Value string `xml:"uuid"`
}
