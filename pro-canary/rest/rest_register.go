package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"sync"
	"time"
)

type Module struct {
	Source     string        //来源类型: Source : scm
	Method     int           //处理方式: Method : 1 build ,0 unbuild
	ModuleName string        //模块名称: ModuleName : nginx-1.1
	DeployPath string        //部署位置: DeployPath : /
	CreateTime time.Duration //创建时间: CreateTime :
	UseVersion int64         //使用版本: UseVersion  : 1
	NewVersion int64         //最新版本: NewVersion : 3
	HaveTar    bool          //是否压缩: HaveTar : true
	Exec       string        //启动命令: Exec : bin/control start
	ExcludeDir []string      //过滤目录: ExcludeDir : logs,data
	Depend     []string      //依赖服务: Depend : mysql,php
	ConfDir    string        //配置路径: ConfDir : conf
}

var store = map[string]*Module{}
var lock = sync.RWMutex{}

//注册软件
func getRegisterModuleTask(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	countries := make([]Module, len(store))
	i := 0
	for _, module := range store {
		countries[i] = *module
		i++
	}
	lock.RUnlock()
	w.WriteJson(&countries)
}

//获取注册的软件
func postRegisterModuleTask(w rest.ResponseWriter, r *rest.Request) {
	module := Module{}
	err := r.DecodeJsonPayload(&module)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if module.ModuleName == "" {
		rest.Error(w, "module code required", 400)
		return
	}
	if module.DeployPath== "" {
		rest.Error(w, "module name required", 400)
		return
	}
	lock.Lock()
	store[module.ModuleName] = &module
	lock.Unlock()
	w.WriteJson(&module)
}
