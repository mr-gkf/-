package controllers

import (
	"baselinCheck/enums"
	"baselinCheck/models"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

//TemplateController 模板管理
type TemplateController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *TemplateController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid", "UpdateSeq")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

//Index 模板管理首页
func (c *TemplateController) Index() {
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "template/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "template/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("TemplateController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("TemplateController", "Delete")
}

// DataGrid 模板管理首页 表格获取数据
func (c *TemplateController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.TemplateQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	data, total := models.TemplatePageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

//Edit 添加、编辑自定义模板
func (c *TemplateController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := models.Template{Id: Id}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["m"] = m
	c.setTpl("template/edit.html", "shared/layout_page.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "template/edit_headcssjs.html"
	c.LayoutSections["footerjs"] = "template/edit_footerjs.html"
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor("TemplateController.Index")
}

//Save 添加、编辑页面 保存
func (c *TemplateController) Save() {
	var err error
	m := models.Template{}
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "提交表单数据失败，可能原因："+err.Error(), m.Id)
	}
	o := orm.NewOrm()
	if m.Id == 0 {
		m.Creator = &c.curUser
		if _, err = o.Insert(&m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失败，可能原因："+err.Error(), m.Id)
		}
	} else {
		oM, err := models.TemplateOne(m.Id)
		oM.TemplateName = m.TemplateName
		oM.TypeId = m.TypeId
		oM.ItemIds = m.ItemIds
		if _, err = o.Update(oM); err != nil {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}
	//添加关系
	var relations []models.TemplateItemRel
	result := strings.Split(m.ItemIds, ",")
	for _, ItemId := range result {
		Id, _ := strconv.Atoi(ItemId)
		if Id == 0 {
			continue
		}
		r := models.Item{Id: Id}
		relation := models.TemplateItemRel{Template: &m, Item: &r}
		relations = append(relations, relation)
	}
	if len(relations) > 0 {
		//批量添加
		if _, err := o.InsertMulti(len(relations), relations); err == nil {
			c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "保存失败", m.Id)
		}
	} else {
		c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
	}
}

//Delete 批量删除
func (c *TemplateController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.TemplateBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}

func (c *TemplateController) UpdateSeq() {
	Id, _ := c.GetInt("pk", 0)
	oM, err := models.TemplateOne(Id)
	if err != nil || oM == nil {
		c.jsonResult(enums.JRCodeFailed, "选择的数据无效", 0)
	}
	o := orm.NewOrm()
	if _, err := o.Update(oM); err == nil {
		c.jsonResult(enums.JRCodeSucc, "修改成功", oM.Id)
	} else {
		c.jsonResult(enums.JRCodeFailed, "修改失败", oM.Id)
	}
}

// ItemGrid 漏洞检查项表格获取数据
func (c *TemplateController) ItemGrid() {
	var params = models.ItemQueryParam{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	data, total := models.ItemPageList(&params)
	fmt.Println("查询出来的数据为：", data)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

// 模板管理编辑页 所属类型获取数据
func (c *TemplateController) TypeOption() {
	//直接反序化获取json格式的requestbody里的值
	var params models.TypeQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	data, total := models.TypePageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

// 获取模板管理编辑页面检查策略选项
func (c *TemplateController) TemplateOption() {
	Id := c.GetString(":id")
	//直接反序化获取json格式的requestbody里的值
	var params models.TemplateQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	params.TypeId = Id
	//获取数据列表和总数
	data, total := models.TemplatePageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}
