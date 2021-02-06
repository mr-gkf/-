package models

import "time"

// Result 实体类
type Result struct {
	TaskId                 int       `orm:"pk"` //任务Id
	CheckItemNum           int       //检查项数量
	CheckItemSuccess       int       //通过检查的检查项数目
	CheckItemFalse         int       //未通过检查的检查项数目
	CheckConformProportion string    `orm:"size(32)"` //检查合规率
	ScanEndTime            time.Time `orm:"auto_now_add;type(datetime)"`//扫描结束时间
	ScanResult             string    //检查结果
}

// TableName 设置表名
func (a *Result) TableName() string {
	return ResultTBName()
}
