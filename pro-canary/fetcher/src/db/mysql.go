package db

import (
	. "config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"strconv"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	Success        = 0
	QueryDBFailed  = 1
	ChangeDBFailed = 2
	PreSqlFailed   = 3
	GetIdFailed    = 4
)

type Api struct {
	database string
	DB       *xorm.Engine
}

func (api *Api) InitSchema(table ...interface{}) error {
	if err := api.DB.StoreEngine("InnoDB").Sync2(table...); err != nil {
		return fmt.Errorf("sync database struct error: %v\n", err)
	}
	return nil
}

//初始化数据库
func (api *Api) InitDB(conf DataBaseConfig) {
	var err error

	api.database = conf.DataBase
	sqlUrl := fmt.Sprintf(conf.URL, conf.User, conf.Password, conf.Host, conf.Port, conf.DataBase)
	api.DB, err = xorm.NewEngine("mysql", sqlUrl)
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
	api.DB.ShowSQL = true
	api.DB.ShowDebug = true
	api.DB.ShowErr = true
	api.DB.SetMaxIdleConns(conf.MaxIdleConnections)
	api.DB.SetMaxOpenConns(conf.MaxOpenConnections)
}

// fetcher

func (api *Api) GetNextId(name string) (next int) {
	sql := fmt.Sprintf("select auto_increment from information_schema.`TABLES` where TABLE_SCHEMA='%s' AND TABLE_NAME='%s'", api.database, name)
	results, queryerr := api.DB.Query(sql)
	if queryerr != nil {
		next = 1
	} else {
		next, queryerr = strconv.Atoi(string(results[0]["auto_increment"]))
		if queryerr != nil {
			next = 1
		}
	}
	return
}
