package netapp

import (
	"encoding/xml"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseStruct(t *testing.T) {
	output, err := xml.Marshal(&Base{})
	if err != nil {
		log.Print(err)
	}

	assert.Equal(t, "<netapp></netapp>", string(output))
}
