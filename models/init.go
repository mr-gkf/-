package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// init 初始化
func init() {
	orm.RegisterModel(new(Course), new(BackendUser), new(Resource), new(Role), new(RoleResourceRel), new(RoleBackendUserRel), new(Task), new(Template), new(Item), new(TemplateItemRel), new(Type), new(Result))
}

// TableName 下面是统一的表名管理
func TableName(name string) string {
	prefix := beego.AppConfig.String("db_dt_prefix")
	return prefix + name
}

// BackendUserTBName 获取 BackendUser 对应的表名称
func BackendUserTBName() string {
	return TableName("backend_user")
}

// ResourceTBName 获取 Resource 对应的表名称
func ResourceTBName() string {
	return TableName("resource")
}

// RoleTBName 获取 Role 对应的表名称
func RoleTBName() string {
	return TableName("role")
}

// RoleResourceRelTBName 角色与资源多对多关系表
func RoleResourceRelTBName() string {
	return TableName("role_resource_rel")
}

// RoleBackendUserRelTBName 角色与用户多对多关系表
func RoleBackendUserRelTBName() string {
	return TableName("role_backenduser_rel")
}

func CourseTBName() string {
	return TableName("course")
}

// TemplateTBName 获取 Template 对应的表名称
func TemplateTBName() string {
	return TableName("template")
}

// ItemTBName 获取 Item 对应的表名称
func ItemTBName() string {
	return TableName("item")
}

// TemplateItemRelTBName 模板与漏洞检查项多对多关系表
func TemplateItemRelTBName() string {
	return TableName("template_item_rel")
}

// TypeTBName 获取 Type 对应的表名称
func TypeTBName() string {
	return TableName("type")
}

// TaskTBName 获取 Task 对应的表名称
func TaskTBName() string {
	return TableName("task")
}

// ResultTBName 获取 Result 对应的表名称
func ResultTBName() string {
	return TableName("result")
}
