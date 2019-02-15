package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type s struct {
	MyValue string
}

func TestMakeRequestFromTemplate(t *testing.T) {
	tpl := []byte(`<query><value>{{.MyValue}}</value></query>`)
	e := `<query><value>abc</value></query>`
	r := makeRequestFromTemplate(tpl, s{MyValue: "abc"})
	assert.Equal(t, e, r.String())
}
