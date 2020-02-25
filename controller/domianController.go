package controller

import (
	"baiduDomain/baiduService"
	"baiduDomain/model"
	"encoding/json"
	"errors"
)

type DomainController struct {
	DomainHandler baiduService.DomainHandler
	MainDomain    string //主域名
	SubDomian     string //需要设置的域名
	Nowip         string //需要设置的ip
}

/**
maindomian
subdomian 设置的解析域名
*/
func (domain *DomainController) SetDomianResolve() error {

	resolveListIn := model.ResolveListIn{
		Domain:   domain.MainDomain,
		PageNo:   1,
		PageSize: 100,
	}
	resolveListInJson, err := json.Marshal(resolveListIn)
	if err != nil {
		return err
	}
	resolveListOut := model.ResolveListOut{}
	dr, err := domain.DomainHandler.DomianResolveList(resolveListInJson)
	if err != nil {
		return err
	}
	jsonerr := json.Unmarshal(dr, &resolveListOut)

	if jsonerr != nil {
		return jsonerr
	}
	count := 0
	if resolveListOut.TotalCount > 0 {

		for _, v := range resolveListOut.Result {
			if v.Domain == domain.SubDomian {
				if domain.Nowip != v.Rdata {
					count++
					//修改解析记录
					resolveEditIn := model.ResolveEditIn{
						v.Domain,
						v.Rdtype,
						v.View,
						domain.Nowip,
						v.TTL,
						v.ZoneName,
						v.RecordID,
					}
					v.Rdata = domain.Nowip
					resolveListEditJson, err := json.Marshal(resolveEditIn)
					if err != nil {
						return err
					}
					re, err := domain.DomainHandler.DomianResolveEdit(resolveListEditJson)
					if err != nil {
						return err
					}
					if !re {
						return errors.New("修改失败")
					}
				} else {
					return errors.New("不需要修改")
				}
			}
		}
	}
	if count == 0 {
		//添加解析记录
		//修改解析记录
		resolveAddIn := model.ResolveEditIn{
			domain.SubDomian,
			"A",
			"DEFAULT",
			domain.Nowip,
			300,
			domain.MainDomain,
			1,
		}
		resolveListAddJson, err := json.Marshal(resolveAddIn)
		if err != nil {
			return err
		}
		re, err := domain.DomainHandler.DomianResolveAdd(resolveListAddJson)
		if err != nil {
			return err
		}
		if !re {
			return errors.New("新增失败")
		}
	}
	return nil
}
