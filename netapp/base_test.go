package netapp

import (
	"encoding/xml"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseStruct(t *testing.T) {
	output, err := xml.Marshal(&Base{Version: "1.7"})
	if err != nil {
		log.Print(err)
	}

	assert.Equal(t, "<netapp version=\"1.7\"></netapp>", string(output))
}
