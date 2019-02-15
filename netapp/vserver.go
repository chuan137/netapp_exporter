package netapp

type Vserver struct {
	Base
	DesiredAttribtues *VserverInfo `xml:"desired-attributes>vserver-info"`
}

type VserverInfo struct {
	VserverName string `xml:"vserver-name,omitempty"`
	UUID        string `xml:"uuid,omitempty"`
}
