package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// TableName 设置Template表名
func (a *Template) TableName() string {
	return TemplateTBName()
}

// TemplateQueryParam 用于搜索的类
type TemplateQueryParam struct {
	BaseQueryParam
	TemplateNameLike string
	TypeId           string
}

// Template 模板实体类
type Template struct {
	Id              int //模板Id
	TypeId          string
	TemplateName    string             `orm:"size(20)"`                    //模板名称
	ItemIds         string             `orm:"size(2000)"`                  //所有已选中的检查项的Id
	TemplateItemRel []*TemplateItemRel `orm:"reverse(many)"`               // 设置一对多的反向关系
	Creator         *BackendUser       `orm:"rel(fk)"`                     //设置一对多关系
	Created         time.Time          `orm:"auto_now_add;type(datetime)"` //创建时间
}

// TemplatePageList 获取分页数据
func TemplatePageList(params *TemplateQueryParam) ([]*Template, int64) {
	query := orm.NewOrm().QueryTable(TemplateTBName())
	data := make([]*Template, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query = query.Filter("templatename__istartswith", params.TemplateNameLike)
	if params.TypeId != "" {
		query = query.Filter("typeid", params.TypeId) // WHERE id = 1
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).RelatedSel().All(&data)
	return data, total
}

// TemplateBatchDelete 批量删除
func TemplateBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(TemplateTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// TemplateOne 获取单条
func TemplateOne(id int) (*Template, error) {
	o := orm.NewOrm()
	m := Template{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
