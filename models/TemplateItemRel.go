package models

import "time"

//建立模板检查项关系表
type TemplateItemRel struct {
	Id         int
	Template   *Template `orm:"rel(fk)"` // 外键
	Item       *Item     `orm:"rel(fk)"` // 外键
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
}

// TableName 设置表名
func (a *TemplateItemRel) TableName() string {
	return TemplateItemRelTBName()
}
