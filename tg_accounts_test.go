package main

import "testing"

func Test_readConf(t *testing.T) {
	gf := readConf()
	addUrlParams(gf)
	err := saveConf(gf)
	if err != nil {
		t.Fatal(err)
	}
}
