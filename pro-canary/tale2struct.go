package main

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"relect"
)

type ModuleInfo struct {
	Source      string   //来源类型: Source : scm
	Method      int      //处理方式: Method : 1 build ,0 unbuild
	ModuleName  string   //模块名称: ModuleName : nginx-1.1
	DeployPath  string   //部署位置: DeployPath : /
	ModuleType  bool     //是否压缩: ModuleType : true
	Exec        string   //启动命令: Exec : bin/control start
	ConfDir     string   //配置路径: ConfDir : conf
	ExcludeDir  []string //过滤目录: ExcludeDir : logs,data
	Depend      []string //依赖服务: Depend : mysql,php
	Description string   //描述: Description : text
}

func fetchAllRowsPtr(query string, struc interface{}, cond ...interface{}) *[]interface{} {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "root", "12345", "127.0.0.1", "3306", "orp")
	db, _ := sql.Open("mysql", conn)
	defer db.Close()
	stmt, err := db.Prepare("SELECT * FROM ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(cond...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	result := make([]interface{}, 0)
	s := reflect.ValueOf(struc).Elem()
	leng := s.NumField()
	onerow := make([]interface{}, leng)
	for i := 0; i < leng; i++ {
		onerow[i] = s.Field(i).Addr().Interface()
	}
	for rows.Next() {
		err = rows.Scan(onerow...)
		if err != nil {
			panic(err)
		}
		result = append(result, s.Interface())
	}
	return &result
}

func main() {
	R := fetchAllRowsPtr("select * from ?", ModuleInfo, "module")
	fmt.Sprintf("%s", R)
}
