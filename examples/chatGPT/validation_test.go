package main

import (
	"io/ioutil"
	"testing"

	"cuelang.org/go/cue/cuecontext"
)

func TestValidation(t *testing.T) {
	host := `Host: "http://mytest:1234"`
	constraints, err := ioutil.ReadFile("constraints.cue")
	if err != nil {
		t.Fatal(err)
	}
	configuration, err := ioutil.ReadFile("wellknown.cue")
	if err != nil {
		t.Fatal(err)
	}
	content := append(constraints, configuration...)
	content = append([]byte(host+"\n"), content...)
	ctx := cuecontext.New()
	v := ctx.CompileBytes(content)
	v = v.Lookup("configuration")
	var aiplugin AIPlugin
	err = v.Decode(&aiplugin)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(aiplugin)
}
