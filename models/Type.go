package models

import (
	"github.com/astaxie/beego/orm"
)

//漏洞检查分类表
type Type struct {
	TypeId             int    `orm:"pk;auto"`
	PropertyType       string `orm:"size(50)"`  //资产类别
	PropertyTypeDetail string `orm:"size(200)"` //资产类别
}

// TypeQueryParam 用于搜索的类
type TypeQueryParam struct {
	BaseQueryParam
}

// TypeName 设置Template表名
func (a *Type) TableName() string {
	return TypeTBName()
}

// TypePageList 获取分页数据
func TypePageList(params *TypeQueryParam) ([]*Type, int64) {
	query := orm.NewOrm().QueryTable(TypeTBName())
	data := make([]*Type, 0)
	//默认排序
	sortorder := "TypeId"
	switch params.Sort {
	case "TypeId":
		sortorder = "TypeId"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}
