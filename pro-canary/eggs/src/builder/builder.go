package main

import (
	"fmt"
	"time"
)

//build分为prebuild 进行一系列的规整工作, 将deploydir打包, 去除conf文件,一些顺序的组合等.
//sourcebuild 进行源代码的编译,属于后期工作.
//目标:生成按模块为单位,线上直接解压的安全可用的包.

//将fetch过来的文件进行预处理,处理成线上可以直接使用的模块包.
//存储prebuild后的内容

//保存的DB
type ModuleBuild struct {
	Id        int
	Category  string `xorm:"notnull unique(name)"` //分类 conf code data
	OriginId  int    `xorm:"notnull unique(name)"` //前两项等同于唯一索引，可以唯一确定模块
	Version   string
	LocalPath string
	Md5String string
	CreatedAt time.Time `xorm:"created"`
	TouchedAt time.Time `xorm:"updated"`
}

type Builder struct {
	BuilerDir  string
	FetcherDir string
	DB         Api.DB
	Md5String  string
}

func NewBuilder() *Builder {
	//加载配置
	configPath := "eggs.toml"
	conf := ConfigFromFile(configPath)
	api := Api{}
	return &Builder{BuilderDir: conf.Packer.BuilderDir, FetcherDir: conf.Packer.FetcherDir, DB: api.DB}
}

func Build(data *ModuleInfo) error {
	dbBuild := ModuleBuild{}
	builder := NewBuilder()
	if err := builder.PreBuild(data); err != nil {
		return err
	}
	//插入数据到数据库
	dbBuild = ModuleBuild{
		Category:  data.Category,
		OriginId:  data.OriginId,
		Version:   data.Version,
		LocalPath: data.SavaPath + "/" + data.FileName,
		Md5String: builder.Md5String,
	}
	if _, err := builder.DB.Insert(&dbBuild); err != nil {
		return err
	}
	return
}

func (b *Builder) PreBuild(data *ModuleInfo) error {
	//创建目录
	fullpath := fmt.Sprintf("%s/%s", b.BaseDir, data.SavaPath)
	_, err := PathExist(fullpath)
	if err == nil {
		err.Error("have the same name.")
		return err
	}
	if err := os.MkdirAll(fullpath, os.ModePerm); err != nil {
		return err
	}

	//判断fetcher目录是否存在.
	fetcherPackage := fmt.Sprintf("%s/%s/%s.tar.gz", b.FetcherDir, data.SavaPath, data.FileName)
	_, err = PathExist(fullpath)
	if err != nil {
		return err
	}

	//解压文件压缩文件,并且会自动创建deploy目录
	deploypath := fmt.Sprintf("%s/%s", fullpath, data.DeployPath)
	if err := UnTarGz(fetcherPackage, deploypath); err != nil {
		return err
	}

	tarFilename := fmt.Sprintf("%s.tar.gz", data.FileName)
	//打包暂时不进行配置与顺序的处理,默认来源是工整的.scm默认有output目录.
	if err := MvDirFiles(deploypath+"/output/", deploypath); err != nil {
		return err
	}
	if err := os.Remove(deploypath + "/output"); err != nil {
		return err
	}

	if err := TarGz(fullpath, tarFileName); err != nil {
		return err
	}
	//生成md5
	md5string, err := MakeMd5(fullpath+"/"+tarFilename, true)
	if err != nil {
		return err
	}
	b.Md5String = md5string

}
