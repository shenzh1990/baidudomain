package main

import (
	"baiduDomain/baiduService"
	"baiduDomain/controller"
	"baiduDomain/utils"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"time"
)

func main() {
	domain()
	go exec()
	select {}
}

func exec() {
	gocron.Every(720).Seconds().Do(domain)
	<-gocron.Start()
}
func domain() {
	handler := baiduService.DomainHandler{
		"test",
		"test",
		"https://bcd.baidubce.com",
		"v1",
	}
	nowip, err := utils.GetIp("4")
	domainController := controller.DomainController{
		handler,
		"qinyule.com",
		"nextcloud",
		nowip,
	}
	timespan := time.Now().String()
	if err != nil {
		fmt.Println(timespan + ":error-" + err.Error())
	} else {
		err := domainController.SetDomianResolve()
		if err != nil {
			fmt.Println(timespan + ":error-" + err.Error())
		} else {
			fmt.Println(timespan + ":success-" + domainController.SubDomian + "设置为" + nowip)
		}
	}

}
