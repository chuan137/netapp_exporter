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
)

func init() {
	url = os.Getenv("NETAPP_URL")
	username = os.Getenv("NETAPP_USERNAME")
	password = os.Getenv("NETAPP_PASSWORD")
}

func TestVserver(t *testing.T) {
	// v := &Vserver{Base: Base{Version: "1.7"}}
	v := &Vserver{}
	v.DesiredAttribtues = &VserverInfo{VserverName: "1"}

	r, e := xml.MarshalIndent(v, "", "\t")
	if e != nil {
		log.Fatal(e)
	}

	exp := `<netapp>
	<desired-attributes>
		<vserver-info>
			<vserver-name>1</vserver-name>
		</vserver-info>
	</desired-attributes>
</netapp>`
	assert.Equal(t, exp, string(r))
}

func TestVserverGet(t *testing.T) {
	c := &Client{url: url, username: username, password: password}
	b := &Base{Client: *c}
	v := &Vserver{Base: *b}

	assert.Equal(t, v.Base.Client.url, "abc")
}
