package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"time"
//	"fmt"
)

//curl -i  -H 'Content-Type: application/json' -H 'Accept: application/json' -u orp:orp -X POST -d {"ModuleList":[{"ModuleName":"nginx","ModuleVersion":1},{"ModuleName":"mysql","ModuleVersion":2}]}' http://127.0.0.1:8080/packer
type Packer struct {
	Id         int64     `json:"id"` //task id
	ModuleList []List    //需要packer的包列表
	CreatedAt  time.Time `xorm:"created"`
}

//需求的包名与包版本
type List struct {
	ModuleName     string `xorm:"uniqe"`
	ModuleVersion int64  `xorm:"version"`
}

func (api *Api) GetPacker(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	packer := Packer{}
	if _, err := api.DB.Where("id=?", id).Get(&packer); err != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&packer)
}

func (api *Api) PostPacker(w rest.ResponseWriter, r *rest.Request) {
	packer := Packer{}
	if err := r.DecodeJsonPayload(&packer); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := api.DB.Insert(&packer); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&packer)
}
