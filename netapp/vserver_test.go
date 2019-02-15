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
	base = Base{Version: "1.7"}
}

func TestVserverDesiredAttributes(t *testing.T) {
	type attr struct {
		VserverName
		UUID
	}

	// set attr to the DesiredAttributes field
	v := &Vserver{base, attr{}}

	r, e := xml.MarshalIndent(v, "", "\t")
	if e != nil {
		log.Fatal(e)
	}

	exp := `<netapp version="1.7">
	<desired-attributes>
		<vserver-info>
			<vserver-name></vserver-name>
			<uuid></uuid>
		</vserver-info>
	</desired-attributes>
</netapp>`
	assert.Equal(t, exp, string(r))
}

func TestVserverGet(t *testing.T) {
	v := &Vserver{Base: base}
	v.Base.Client = &Client{config: &ClientConfig{url, username, password}}

	assert.Equal(t, v.Base.Client.config.url, "abc")
}
