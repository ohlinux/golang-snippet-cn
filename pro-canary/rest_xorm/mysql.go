package main

import (
	//	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
//	"os"
	//	"sync"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	DataBase   = "orp"
	configPath = "eggs.toml"
)

const (
	Success        = 0
	QueryDBFailed  = 1
	ChangeDBFailed = 2
	PreSqlFailed   = 3
	GetIdFailed    = 4
)

type Api struct {
	DB *xorm.Engine
}

func (api *Api) InitSchema(table interface{}) error {
    if err := api.DB.StoreEngine("InnoDB").Sync2(table); err != nil {
		return fmt.Errorf("sync database struct error: %v\n", err)
	}
        return  nil
 }

//初始化数据库
func (api *Api) InitDB(database string) {
	var err error
	conf, err := ConfigFromFile(configPath)
	if err != nil {
		panic(err.Error())
	}

	dbConf := conf.DataBase
	sqlUrl := fmt.Sprintf(dbConf.URL, dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, database)
	api.DB, err = xorm.NewEngine("mysql", sqlUrl)
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
	api.DB.ShowSQL = true
	api.DB.ShowDebug = true
	api.DB.ShowErr = true
	api.DB.SetMaxIdleConns(dbConf.MaxIdleConnections)
	api.DB.SetMaxOpenConns(dbConf.MaxOpenConnections)

//	f, err := os.Create("sql.log")
//	if err != nil {
//		println(err.Error())
//		return
//	}
//	defer f.Close()
//	api.DB.Logger = xorm.NewSimpleLogger(f)

}
