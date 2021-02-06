package controllers

import (
	"baselinCheck/enums"
	"baselinCheck/models"
	"baselinCheck/utils"
	"encoding/json"
	"time"

	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

//TaskController 课程管理
type TaskController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *TaskController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid", "DataList")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

//Index 任务管理首页
func (c *TaskController) Index() {
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "task/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "task/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("TaskController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("TaskController", "Delete")
}

// DataGrid 任务管理首页 表格获取数据
func (c *TaskController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.TaskQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	data, total := models.TaskPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

//Edit 添加、编辑任务界面
func (c *TaskController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := models.Task{Id: Id, StartTime: time.Now()}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["m"] = m
	c.setTpl("task/edit.html", "shared/layout_page.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "task/edit_headcssjs.html"
	c.LayoutSections["footerjs"] = "task/edit_footerjs.html"
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor("TaskController.Index")
}

//Save 添加、编辑页面 保存
func (c *TaskController) Save() {
	var err error
	m := models.Task{}
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "提交表单数据失败，可能原因："+err.Error(), m.Id)
	}
	o := orm.NewOrm()
	if m.Id == 0 {
		Template := new(models.Template)
		Template.Id, _ = strconv.Atoi(m.TemplateRelId)
		m.Template = Template
		Type := new(models.Type)
		Type.TypeId, _ = strconv.Atoi(m.TypeRelId)
		m.Type = Type
		m.Password = utils.String2md5(m.Password)
		m.SuperPassword = utils.String2md5(m.SuperPassword)
		m.Creator = &c.curUser
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失败，可能原因："+err.Error(), m.Id)
		}
	} else {
		oM, err := models.TaskOne(m.Id)
		oM.Name = m.Name
		oM.ScanTarget = m.ScanTarget
		oM.ExeType = m.ExeType
		Template := new(models.Template)
		Template.Id, _ = strconv.Atoi(m.TemplateRelId)
		oM.Template = Template
		oM.TemplateRelId = m.TemplateRelId
		Type := new(models.Type)
		Type.TypeId, _ = strconv.Atoi(m.TypeRelId)
		oM.Type = Type
		oM.TypeRelId = m.TypeRelId
		fmt.Println("oM.TypeRelId为：", oM.TypeRelId)
		oM.Account = m.Account
		oM.ConMethod = m.ConMethod
		oM.StartTime = m.StartTime
		oM.Password = strings.TrimSpace(m.Password)
		if len(oM.Password) == 0 {
			//如果密码为空则不修改
			oM.Password = m.Password
		} else {
			oM.Password = utils.String2md5(m.Password)
		}
		oM.Port = m.Port
		oM.SuperAccount = m.SuperAccount
		if len(oM.SuperPassword) == 0 {
			//如果密码为空则不修改
			oM.SuperPassword = m.SuperPassword
		} else {
			oM.SuperPassword = utils.String2md5(m.SuperPassword)
		}
		if _, err = o.Update(oM); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			utils.LogDebug(err)
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}
}

//Delete 批量删除
func (c *TaskController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.TaskBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}

// 执行任务
func (c *TaskController) ExecuteTask() {
	Id, _ := c.GetInt(":id", 0)
	o := orm.NewOrm()
	oM, err := models.TaskOne(Id)
	if err != nil || oM == nil {
		c.jsonResult(enums.JRCodeFailed, "选择的数据无效", 0)
	}
	oM.IsExcute = 1
	if _, err = o.Update(oM); err == nil {
		fmt.Println("更新数据成功")
	} else {
		fmt.Println("更新数据失败", err)
	}
	t := &utils.Task{}
	t.TaskID = oM.Id
	t.Ip = oM.ScanTarget
	switch oM.ConMethod {
	case 1:
		t.Protocol = "ssh"
	case 2:
		t.Protocol = "smb"
	default:
		t.Protocol = "telent"
	}
	t.Port = oM.Port
	t.UserName = oM.Account
	t.PassWord = oM.Password
	t.SuperUserName = oM.SuperAccount
	t.SuperPassWord = oM.SuperPassword
	t.PropertyType = oM.Type.PropertyType
	ids := make([]string, 0)
	for _, str := range strings.Split(oM.Template.ItemIds, ",") {
		ids = append(ids, str)
	}
	var params models.ItemQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	params.Ids = ids
	M, num := models.ItemList(&params)
	t.Num = num
	for key, _ := range M {
		t.Item = append(t.Item, utils.Item{CheckItemId: M[key].ItemId})
	}
	//定义返回的数据结构
	result := make(map[string]interface{})
	if err != nil {
		result["code"] = 2
		result["msg"] = "任务执行失败"
		c.Data["json"] = result
		c.ServeJSON()
	} else {
		result["code"] = 1
		result["msg"] = "任务已经执行"
		c.Data["json"] = result
		c.ServeJSON()
	}
	utils.SendMessage(*t)
	go func() {
		for {
			scanResult := utils.GetResult()
			saveResult(scanResult)
		}
	}()
}

// 扫描结果保存入库
func saveResult(result *utils.Result) {
	var err error
	m := models.Result{}
	o := orm.NewOrm()
	m.TaskId = result.TaskID
	m.CheckConformProportion = result.CheckConformProportion
	m.CheckItemFalse = result.CheckItemFalse
	m.CheckItemSuccess = result.CheckItemSuccess
	m.CheckItemNum = result.CheckItemNum
	var detail string
 	for _, value := range result.Item {
		detail += value.CheckItemId + "," + value.CheckItemInfo + "," + value.CheckItemResult + ";"
	}
	m.ScanResult = detail
	if _, err = o.Insert(&m); err == nil {
		fmt.Println("插入数据库成功")
	} else {
		fmt.Println("插入数据库失败", err)
	}
}