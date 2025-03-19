package main

var accountMap map[string]string // key是username,value是pwd

func init() {
	InitAccounts()

	gf := readConf()
	editConf(gf)
	err := saveConf(gf)
	if err != nil {
		panic(err)
	}

	accountMap = make(map[string]string, len(gf))
	for _, c := range gf {
		if _, ok := accountMap[c.Username]; !ok {
			accountMap[c.Username] = c.Pwd
		}
	}

	// 配置糖果:8888代理
	writeTgProxyConf()
}
