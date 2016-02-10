package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"strconv"
	"time"
)

//curl -i  -H 'Content-Type: application/json' -H 'Accept: application/json' -u orp:orp -X POST -d '{"MenuId": 1 , "AppId": 2 , "CallBack": "http://xxx.com/yyy","Modules":[{"OriginId":1,"Version":"1.1.1","Category":"bin"},{"OriginId":3,"Version":"1.2.1","Category":"conf"}]}' http://127.0.0.1:8080/packer/module

//curl -i  -H 'Content-Type: application/json' -H 'Accept: application/json' -u orp:orp -X POST -d '{"MenuId": 1 , "AppId": 2 , "CallBack": "http://xxx.com/yyy"}' http://127.0.0.1:8080/packer/app

/*
request data
modules : [ {"id": 1 , "version": "1.1.1","type":"bin"} ]
callback: url
menuid : ?
AppId:
*/

//request module data
type PackerModule struct {
	Id        int
	MenuId    int //menuid上线单id跟
	AppId     int
	Modules   []RequestModule //需要packer的包列表
	CallBack  string          //回调地址
	Status    int             //处理过程中的状态.
	CreatedAt time.Time       `xorm:"created"`
}

//需求的包名与包版本还有类型
type RequestModule struct {
	OriginId int //模块id 或者conf id
	Version  string
	Category string //类型:conf bin data
}

//request app data
type PackerApp struct {
	Id       int
	MenuId   int
	AppId    int
	CallBack string
	Status   int
	CreateAt time.Time
}

//callback json
type CallBack struct {
	ErrNo  int // 0 成功 1失败
	ErrMsg string
	Data   CallBackData
}

//callback return data
type CallBackData struct {
	MenuId int
	AppId  int
	URL    string
	Md5    string
}

//table app_program_ori
type AppProgramOri struct {
	Id         int    `xorm:"int(11)"`
	Type       string `xorm:"varchar(20)"`
	SrcPath    string `xorm:"varchar(200) notnull"`
	Uncompress int    `xorm:"int(10)"` //是否解压缩，1代表需要，-1代表不需要解压缩
	DeployPath string `xorm:"varchar(45) notnull"`
}

func (api *Api) GetPacker(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	packer := PackerModule{}
	has, err := api.DB.Where("id=?", id).Get(&packer)
	if err != nil {
		rest.NotFound(w, r)
		return
	}
	if !has {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&packer)
}

func (api *Api) PostPackerModule(w rest.ResponseWriter, r *rest.Request) {
	packer := PackerModule{}
	dbProgram := AppProgramOri{}

	if err := r.DecodeJsonPayload(&packer); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//将任务插入到表格中.
	if _, err := api.DB.Insert(&packer); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//数据重组,减少任务传输中的查询 .
	modules := []ModuleInfo{}

	for _, va := range packer.Modules {

		//获取必要的信息
		_, err := api.DB.Where("id=?", va.OriginId).Get(&dbProgram)
		if err != nil {
			rest.NotFound(w, r)
			return
		}

		uniqstring := fmt.Sprintf("%d%s%s", va.OriginId, va.Category, va.Version)

		moduleinfo := ModuleInfo{
			AppId:      packer.AppId,
			OriginId:   va.OriginId,
			Version:    va.Version,
			Category:   va.Category,
			SrcPath:    dbProgram.SrcPath,
			SrcType:    dbProgram.Type,
			Compress:   dbProgram.Uncompress,
			DeployPath: dbProgram.DeployPath,
			FileName:   MakeFileName(dbProgram.SrcPath, packer.AppId, va.OriginId),
			SavaPath:   MakePathName(uniqstring),
		}
		modules = append(modules, moduleinfo)
		fmt.Println(moduleinfo)
	}

	//分配任务给各个worker

	api.DP.Exec(modules, TaskData(strconv.Itoa(packer.MenuId)))

	w.WriteJson(&packer)
}

func (api *Api) PostPackerApp(w rest.ResponseWriter, r *rest.Request) {
	packer := PackerApp{}
	if err := r.DecodeJsonPayload(&packer); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := api.DB.Insert(&packer); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//分配任务前先获取任务数据
	//api.DP.Exec(packer.Modules, TaskData(strconv.FormatInt(packer.Id, 10)))
	w.WriteJson(&packer)
}
