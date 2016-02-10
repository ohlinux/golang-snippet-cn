/*
usage : curl -i -u orp:orp http://127.0.0.1:8080/register/module
*/
package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"os"
)

func LaunchRestServer() {

	//加载rest配置
	//config, err := ConfigFromFile(configPath)
	//if err != nil {
	//	panic(err.Error())
	//}

	//初始化api
	api := Api{}
	api.InitDB("orp")
	//更新数据库表,默认不会删除多余的字段.
	api.InitSchema(new(Module))
	api.InitSchema(new(Packer))
	//初始化任务调度器
	api.InitDispatcher()
	//启动任务调度器,有三个worker在监听服务.
	go api.DP.Start()

	var logErr error
	restLog, logErr := os.OpenFile("./log/rest_access.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if logErr != nil {
		panic(logErr)
	}
	defer restLog.Close()
	restErrLog, logErr := os.OpenFile("./log/rest_error.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if logErr != nil {
		panic(logErr)
	}
	defer restErrLog.Close()

	handler := rest.ResourceHandler{
		PreRoutingMiddlewares: []rest.Middleware{
			&rest.AuthBasicMiddleware{
				Realm: "eggs zone",
				//AllowedMethods: []string{"GET", "POST", "PUT"},
				Authenticator: func(userId string, password string) bool {
					if userId == "orp" && password == "orp" {
						return true
					}
					return false
				},
			},
		},
		EnableRelaxedContentType: true,
		Logger:      log.New(restLog, "", 0),
		ErrorLogger: log.New(restErrLog, "", log.Ldate|log.Ltime|log.Llongfile),
	}

	//设置路由
	err := handler.SetRoutes(
		//register软件绑定
		&rest.Route{"GET", "/register/module", api.GetAllModules},
		&rest.Route{"POST", "/register/module", api.PostModule},
		&rest.Route{"GET", "/register/module/:id", api.GetModule},
		&rest.Route{"PUT", "/register/module/:id", api.PutModule},
		&rest.Route{"DELETE", "/register/module/:id", api.DeleteModule},

		//packer打包需求
		&rest.Route{"POST", "/packer", api.PostPacker},
		&rest.Route{"GET", "/packer/:id", api.GetPacker},
		//version版本查询
		//&rest.Route{"GET", "/version/diff/:vnew/:vold", getVersionDiffTask},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":8080", &handler))

	api.DP.Stop()
}

//func main() {
//	LaunchRestServer()
//}
