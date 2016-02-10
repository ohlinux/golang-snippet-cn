package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"sync"
	"time"
	//"fmt"
)

//example
//curl -i  -H 'Content-Type: application/json' -H 'Accept: application/json' -X POST -u orp:orp -d  '{"Source":"scm","Method":0,"ModuleName":"nginx-1.1","DeployPath":"/","CreateTime":1418979835,"UseVersion":11,"NewVersion":15,"HaveTar":true,"Exec":"bin/control start","ExcludeDir":["logs","data"],"ConfDir":"conf","Depend":["mysql","php"]}'  http://127.0.0.1:8080/register/module
type Module struct {
	Id          int64  `json:"id"`                                    //module id
	AppName     string `xorm:"size:255 notnull`                       //关联的app名字
	ModuleName  string `xorm:"size:255 notnull unique 'module_name'"` //模块名称: 同一个app的modulename不可重复.
	SourceType  string //来源类型: Source : scm
	Compress    bool   //处理方式: true压缩了,需要解压,false 不需要解压.
	DeployPath  string //部署位置: DeployPath : /
	Exec        string //启动命令: Exec : bin/control start
	ConfDir     string //配置路径: ConfDir : conf
	ExcludeDir  string //忽略目录: ExcludeDir : logs,data
	Depend      string //依赖服务: Depend : mysql,php
	Description string //描述: Description : text
	UseVersion  string //版本号统一用string进行存储
	LastVersion string
	CreatedAt   time.Time `xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated"`
	DeletedAt   time.Time `xorm:"deleted"`
}

var store = map[string]*Module{}

func (api *Api) GetAllModules(w rest.ResponseWriter, r *rest.Request) {
	modules := []Module{}
	api.DB.Find(&modules)
	w.WriteJson(&modules)
}

func (api *Api) GetModule(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	module := Module{}
	//result,err := api.DB.Where("id=?", id).Get(&module){
	//api.DB.Id(id).Get(&module)
	if _, err := api.DB.Where("id=?", id).Get(&module); err != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&module)
}

func (api *Api) PostModule(w rest.ResponseWriter, r *rest.Request) {
	module := Module{}
	if err := r.DecodeJsonPayload(&module); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := api.DB.Insert(&module); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&module)
}

func (api *Api) PutModule(w rest.ResponseWriter, r *rest.Request) {

	id := r.PathParam("id")
	module := Module{}
	if _, err := api.DB.Id(id).Update(&module); err != nil {
		rest.NotFound(w, r)
		return
	}

	updated := Module{}
	if err := r.DecodeJsonPayload(&updated); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	module.ModuleName = updated.ModuleName

	if _, err := api.DB.Id(id).Update(&module); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&module)
}

func (api *Api) DeleteModule(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	module := Module{}
	if _, err := api.DB.Id(id).Delete(&module); err != nil {
		rest.NotFound(w, r)
		return
	}
	if _, err := api.DB.Id(id).Delete(&module); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
