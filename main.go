package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Department struct {
	Id   string
	Name string
}

var kf = Department{
	Id:   "1",
	Name: "开发部",
}
var sc = Department{
	Id:   "2",
	Name: "市场部",
}

type Employee struct {
	Id    string
	Name  string
	DepId string
}

var xm = Employee{
	Id:    "1",
	Name:  "小明",
	DepId: "1",
}

var xh = Employee{
	Id:    "2",
	Name:  "小红",
	DepId: "2",
}

// 用户表
// 部门表(域表)
// 角色表
// 如需要 路径 修改为 功能
// act 修改为 具体动作


func main() {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/am"))
	if err != nil {
		log.Fatal(err)
	}

	// 自动保存 gorm
	a, err := gormadapter.NewAdapterByDB(db) // Your driver and data source.
	if err != nil {
		log.Fatal(err)
	}

	// casbin client
	enforcer, err := casbin.NewEnforcer("./model.conf", a)
	if err != nil {
		log.Fatalf("error, detail: %s", err)
	}

	SavePolicy(enforcer)
	Validate(enforcer)
	Show(enforcer)

}

func SavePolicy(enforcer *casbin.Enforcer) {
	// 1 为 开发部门 角色1
	_, err := enforcer.AddPolicy("1", kf.Id, "/admin", "post")
	if err != nil {
		log.Fatal(err)
	}
	_, err = enforcer.AddPolicy("1", kf.Id, "/admin", "delete")
	if err != nil {
		log.Fatal(err)
	}
	// 2 为 开发部门 角色2
	_, err = enforcer.AddPolicy("2", kf.Id, "/dep", "post")
	if err != nil {
		log.Fatal(err)
	}
	_, err = enforcer.AddPolicy("2", kf.Id, "/dep", "delete")
	if err != nil {
		log.Fatal(err)
	}
	// 3 为 市场部 角色3
	_, err = enforcer.AddPolicy("3", kf.Id, "/admin", "post")
	if err != nil {
		log.Fatal(err)
	}
	_, err = enforcer.AddPolicy("3", kf.Id, "/admin", "delete")
	if err != nil {
		log.Fatal(err)
	}

	// 3 为 市场部 角色3
	_, err = enforcer.AddPolicy("3", sc.Id, "/admin", "post")
	if err != nil {
		log.Fatal(err)
	}
	_, err = enforcer.AddPolicy("3", sc.Id, "/admin", "delete")
	if err != nil {
		log.Fatal(err)
	}

	// 4 为 开发部门 角色 4
	_, err = enforcer.AddPolicy("4", sc.Id, "/dep", "post")
	if err != nil {
		log.Fatal(err)
	}
	_, err = enforcer.AddPolicy("4", sc.Id, "/dep", "delete")
	if err != nil {
		log.Fatal(err)
	}

	enforcer.AddRoleForUserInDomain(xm.Id, "1", kf.Id)
	enforcer.AddRoleForUserInDomain(xh.Id, "4", sc.Id)
}

func Validate(enforcer *casbin.Enforcer) {
	requests := [][]interface{}{
		{xm.Id, xm.DepId, "/admin/111", "post"},
		{xm.Id, xm.DepId, "/admin", "post"},
		{xm.Id, xm.DepId, "/dep", "post"},
		{xh.Id, xh.DepId, "/admin/111", "post"},
		{xh.Id, xh.DepId, "/admin", "post"},
		{xh.Id, xh.DepId, "/dep", "post"},
	}
	for i, request := range requests {
		ok, valid, err := enforcer.EnforceEx(request...)
		fmt.Printf("No.%d: valid: %t param: %s err: %v\n", i+1, ok, valid, err)
	}
}

func Show(enforcer *casbin.Enforcer) {
	// 获取 开发部门下的角色列表
	// 虽然写的是Users 用的改conf后 获取的是 角色
	perms := enforcer.GetAllUsersByDomain("1")
	for _, perm := range perms {
		fmt.Println(perm)
	}

	// 获取 开发部 的所有权限
	// 1 表示匹配第二个字段 domain 也就是部门id字段
	polices := enforcer.GetFilteredPolicy(1, kf.Id)
	for _, policy := range polices {
		fmt.Println(policy)
	}
}

