package main

import (
	"fmt"
)

//封装,将modules处理完的结果按一定的规则进行封装. 并且封装完后进行callback的调用告诉前端已经封装完毕.
//将结果记录在任务表格中.PackerModule

type Packer struct {
	PackerDir  string
	BuilderDir string
	DB         Api.DB
	Md5String  string
}

func NewPacker() *Packer {
	//加载配置
	configPath := "eggs.toml"
	conf := ConfigFromFile(configPath)
	api := Api{}
	return &Builder{PackerDir: conf.Packer.PackerDir, BuilderDir: conf.Packer.BuilderDir, DB: api.DB}
}

//打包module
func PackerM(data *Task) error {
	packer := NewPacker()
	//创建存储目录
	filename := "module.tar.gz"
	uniqname := fmt.Sprintf("%d", data.TaskId)
	dir := MakePathName(uniqname)
	fullPath := fmt.Sprintf("%s/%s", packer.PackerDir, dir)

	if err := os.MkdirAll(fullPath+"/"+"files", os.ModePerm); err != nil {
		return err
	}

	//copy文件 tar.gz tar.gz.md5

	for _, va := range data.Data {
		tarFile := fmt.Sprintf("%s/%s/%s.tar.gz", packer.BuilderDir, va.SavaPath, va.FileName)
		dstFile := fmt.Sprintf("%s/%s.tar.gz", fullPath, va.FileName)
		if err := CopyFile(tarFile, dstFile); err != nil {
			return err
		}
		if err := CopyFile(tarFile+".md5", dstFile+"md5"); err != nil {
			return err
		}
	}

	//生成描述文件toml description

}

//打包app
func PackerA() error {

}
