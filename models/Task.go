package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

// TableName 设置Task表名
func (a *Task) TableName() string {
	return TaskTBName()
}

// TaskQueryParam 用于搜索的类
type TaskQueryParam struct {
	BaseQueryParam
	NameLike string
}

// Task 实体类
type Task struct {
	Id            int          //任务Id
	Name          string       `orm:"size(20)"` //任务名称
	Type          *Type        `orm:"rel(fk)"`  //设置一对多关系
	TypeRelId     string       //辅助存储类型下拉选项id用的
	Template      *Template    `orm:"rel(fk)"` //设置一对多关系
	TemplateRelId string       //辅助存储检查策略下拉选项id用的
	ScanTarget    string       `orm:"size(256)"` //扫描目标
	ExeType       int          //执行类型
	StartTime     time.Time    //开始时间
	Account       string       `orm:"size(32)"` //用户名
	Password      string       `json:"-"`       //密码
	ConMethod     int          //连接方式
	Port          int          //端口号
	SuperAccount  string       `orm:"size(32)"` //超级权限登陆账户
	SuperPassword string       `json:"-"`       //超级权限登录密码
	IsExcute      int          //是否执行
	Creator       *BackendUser `orm:"rel(fk)"`                     //设置一对多关系
	Created       time.Time    `orm:"auto_now_add;type(datetime)"` //创建时间
}

// TaskPageList 获取分页数据
func TaskPageList(params *TaskQueryParam) ([]*Task, int64) {
	query := orm.NewOrm().QueryTable(TaskTBName())
	data := make([]*Task, 0)
	//默认排序
	sortorder := "Id"
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query = query.Filter("name__istartswith", params.NameLike)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).RelatedSel().All(&data)
	return data, total
}

// TaskBatchDelete 批量删除
func TaskBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(TaskTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// TaskOne 获取单条
func TaskOne(id int) (*Task, error) {
	o := orm.NewOrm()
	task := &Task{}
	query := o.QueryTable(TaskTBName())
	query = query.Filter("Id", id)
	_, err := query.Count()
	query.RelatedSel().All(task)
	fmt.Println(task.Template)
	return task, err
}
