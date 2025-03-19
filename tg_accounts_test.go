package main

import "testing"

func Test_readConf(t *testing.T) {
	gf := readConf()
	editConf(gf)
	err := saveConf(gf)
	if err != nil {
		t.Fatal(err)
	}
}
