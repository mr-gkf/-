package utils

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
)
type Task struct {
	TaskID        int
	Ip            string
	Protocol      string
	Port          int
	UserName      string
	PassWord      string
	SuperUserName string
	SuperPassWord string
	PropertyType  string
	Num           int
	Item          []interface{}
}

type Item struct {
	CheckItemId string
}

type Result struct {
	TaskID                 int
	CheckItemNum           int
	CheckItemSuccess       int
	CheckItemFalse         int
	CheckConformProportion string
	Num                    int
	Item                   []ScanResult
}

type ScanResult struct {
	CheckItemId     string
	CheckItemResult string
	CheckItemInfo   string
}

func InitRedis() *redis.Client {
	host := beego.AppConfig.String("cache::redis_host")
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       0,
	})
	return rdb
}
// 向Redis中发送数据
func SendMessage(task Task) {
	rdb := InitRedis()
	SendOutputParameter(rdb, task)
}

// 获取扫描结果
func GetResult() *Result {
	m := Result{}
	rdb := InitRedis()
	scanResult, _ := GetInputParameter(rdb, &m)
	return scanResult
}

func SendOutputParameter(client *redis.Client, outputParam Task) (err error) {
	outParamJson, _ := json.Marshal(outputParam)
	err = client.LPush("list1", outParamJson).Err()
	if err != nil {
		fmt.Println("[error]	LPush error: ", err.Error())
		return err
	}
	fmt.Println("[info]	Output parameters: ", outputParam)
	return nil
}

func GetInputParameter(client *redis.Client, inputParam *Result) (scanResult *Result,err error) {
	inputParamJson, err := client.BLPop(0, "list2").Result()
	if err != nil {
		fmt.Println("[error]	BLPop error: ", err.Error())
		return nil, err
	}
	err = json.Unmarshal([]byte(inputParamJson[1]), inputParam)
	if err != nil {
		fmt.Println("[error]	Unmarshal error: ", err.Error())
		return nil, err
	}

	fmt.Println("[info]	Input parameters: ", *inputParam)
	return inputParam, nil
}