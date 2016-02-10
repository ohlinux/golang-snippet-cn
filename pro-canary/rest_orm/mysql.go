package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"sync"
)

var (
	//hostName, user, password, url string
	//port                          int64
	_init_db       sync.Once
	_db_connection *sql.DB
)

func handleError(db *sql.DB) {
	if x := recover(); x != nil {
		fmt.Sprintf("DB connection failed: %v", x)
		db.Close()
	}
}

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

var ErrorCode = map[int]string{
	0: "SUCCESS",
	1: "query db failed",
	2: "change db failed",
	3: "prepare sql failed",
	4: "get insert id failed",
}

type Api struct {
	DB gorm.DB
}

func (api *Api) InitSchema() {
        api.DB.AutoMigrate(&Module{})
    }

func (api *Api) InitDB(database string) {
	var err error
	conf, err := ConfigFromFile(configPath)
	if err != nil {
		panic(err.Error())
	}

	dbConf := conf.DataBase
	sqlUrl := fmt.Sprintf(dbConf.URL, dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, database)
	api.DB, err = gorm.Open("mysql", sqlUrl)
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
	api.DB.DB()
	api.DB.DB().Ping()
	api.DB.DB().SetMaxIdleConns(dbConf.MaxIdleConnections)
	api.DB.DB().SetMaxOpenConns(dbConf.MaxOpenConnections)
	api.DB.LogMode(true)
}

//每次返回的DB是线程安全的，如果多个routine操作一张表，不要复用sql.DB
func GetDbConnection(database string) *sql.DB {
	_init_db.Do(func() {
		conf, err := ConfigFromFile(configPath)
		if err != nil {
			panic(err.Error())
		}

		dbConf := conf.DataBase
		//conn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbuser, dbpass, dbhost, dbport, dbname)
		sqlUrl := fmt.Sprintf(dbConf.URL, dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, database)
		fmt.Println(sqlUrl)
		var dbErr error
		_db_connection, dbErr = sql.Open("mysql", sqlUrl)
		_db_connection.SetMaxIdleConns(dbConf.MaxIdleConnections)
		_db_connection.SetMaxOpenConns(dbConf.MaxOpenConnections)
		if dbErr != nil {
			_db_connection.Close()
			fmt.Sprintf("DB connection failed: %s", dbErr)
			_db_connection = nil
		}
	})
	return _db_connection
}
