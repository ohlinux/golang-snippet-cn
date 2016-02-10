package db

import (
	"time"
)

type Module struct {
	Id          int64     `json:"id"` //module id
	Source      string    //来源类型: Source : scm
	Method      int       //处理方式: Method : 1 build ,0 unbuild
	ModuleName  string    `xorm:"size:255 notnull unique 'module_name'"` //模块名称: ModuleName : nginx-1.1
	DeployPath  string    //部署位置: DeployPath : /
	ModuleType  bool      //是否压缩: ModuleType : true
	Exec        string    //启动命令: Exec : bin/control start
	ConfDir     string    //配置路径: ConfDir : conf
	ExcludeDir  string    //过滤目录: ExcludeDir : logs,data
	Depend      string    //依赖服务: Depend : mysql,php
	Description string    //描述: Description : text
	CreatedAt   time.Time `xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated"`
	DeletedAt   time.Time `xorm:"deleted"`
	UseVersion  string    //使用版本: UseVersion  : 1
	LastVersion int64     //最新版本: LastVersion : 3
}

// 非db struct
type PreFetch struct {
	Type     string
	OriginId int64  // 两项等同于ModuleId，可以唯一确定模块
	AppId    int64  // 用于取配置
	Source   string // 获取方式 scm/http/ftp/配置管理模块
	Name     string // 模块名、路径、…… 用于获取src路径
	Version  string // 要获取的版本
	// Force     bool    // 留空 用于MD5冲突时强制生成新模块
}

type AppOriStatu struct {
	Id        int64
	Type      string
	OriginId  int64 // 两项等同于ModuleId，可以唯一确定模块
	Version   string
	LocalPath string
	Md5String string
	CreatedAt time.Time `xorm:"created"` // log?不懂 看你加了我留着了
	TouchedAt time.Time `xorm:"updated"` // 留给清理陈旧包的工具使用
}
