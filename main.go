package main

import (
	_ "baselinCheck/routers"
	_ "baselinCheck/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
