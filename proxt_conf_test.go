package main

import "testing"

func Test_writeTgProxyConf(t *testing.T) {
	if err := writeTgProxyConf(); err != nil {
		panic(err)
	}
}
