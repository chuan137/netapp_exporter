package netapp

import (
	"encoding/xml"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	url      string
	username string
	password string
	base     Base
)

func init() {
	url = os.Getenv("NETAPP_URL")
	username = os.Getenv("NETAPP_USERNAME")
	password = os.Getenv("NETAPP_PASSWORD")
}

func TestNewVserver(t *testing.T) {
	type query struct {
		VserverType
	}
	type attr struct {
		VserverName
		UUID
	}

	// set attr to the DesiredAttributes field
	v := NewVserver(attr{})
	v.Params.XMLName.Local = "vserver-get-iter"
	v.Params.Query = query{VserverType{"cluster | data"}}

	r, e := xml.MarshalIndent(v, "", "\t")
	if e != nil {
		log.Fatal(e)
	}

	exp := `<netapp version="1.7">
	<vserver-get-iter>
		<query>
			<vserver-info>
				<vserver-type>cluster | data</vserver-type>
			</vserver-info>
		</query>
		<desired-attributes>
			<vserver-info>
				<vserver-name></vserver-name>
				<uuid></uuid>
			</vserver-info>
		</desired-attributes>
	</vserver-get-iter>
</netapp>`
	assert.Equal(t, exp, string(r))
}

func TestVserverGet(t *testing.T) {
	v := &Vserver{Base: base}
	v.Base.Client = &Client{config: &ClientConfig{url, username, password}}

	assert.Equal(t, v.Base.Client.config.url, "abc")
}
