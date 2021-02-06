package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

// Item 漏洞检查项实体类
type Item struct {
	Id               int //主键Id
	ItemId           string
	ItemName         string `orm:"size(100)"` //检查项名称
	Category         int    //检查项所属分类id
	Weight           int    //检查项权重值
	OkDescription    string `orm:"size(4000)"` //检查项合格描述
	FalseDescription string `orm:"size(4000)"` //检查项不合格描述
	EnhancePlan      string `orm:"size(4000)"` //检查项加固方案
}

// ItemQueryParam 用于搜索的类
type ItemQueryParam struct {
	BaseQueryParam
	ItemNameLike string
	CategoryLike string
	Ids          []string
}

// TableName 设置表名
func (a *Item) TableName() string {
	return ItemTBName()
}

func ItemPageList(params *ItemQueryParam) ([]*Item, int64) {
	query := orm.NewOrm().QueryTable(ItemTBName())
	data := make([]*Item, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query = query.Filter("itemname__istartswith", params.ItemNameLike)
	query = query.Filter("category__istartswith", params.CategoryLike)
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	fmt.Println("查询出来的数据为：", data)
	total, _ := query.Count()
	return data, total
}

func ItemList(params *ItemQueryParam) ([]*Item, int) {
	query := orm.NewOrm().QueryTable(ItemTBName())
	data := make([]*Item, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query = query.Filter("id__in", params.Ids)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, int(total)
}
