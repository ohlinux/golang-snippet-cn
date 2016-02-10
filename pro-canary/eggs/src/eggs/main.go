package main

import (
	. "common/config"
	. "common/db"
	"common/logs"
	"common/utils"
	"flag"
	"rest"
)

var configPath string

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//传入配置文件
	//var configPath string
	flag.StringVar(&configPath,
		"conf",
		"conf/config.toml", "Path of the TOML configuration of the eggs.conf ,default conf/config.toml .")
	flag.Parse()

	config, err := ConfigFromFile(configPath)
	if err != nil {
		panic(err.Error())
	}

	//加载
	LoadConf(config.Loging.File)

	//初始化api
	api := Api{}
	api.InitDB("orp")
	//更新数据库表,默认不会删除多余的字段.
	api.InitSchema(new(PackerModule))
	api.InitSchema(new(PackerApp))
	api.InitSchema(new(AppProgramOri), new(ModuleBuild))
	//初始化任务调度器
	api.InitDispatcher()

	//启动任务调度器,有三个worker在监听服务.
	go api.DP.Start()

	//启动rest server
	rest.LaunchRestServer(configPath)

	//reciver signal
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP)
	//#WORKER is a new process tag.
	//newArgs := append(os.Args, "#WORKER")
	attr := syscall.ProcAttr{
		Env: os.Environ(),
	}
	for {
		sig := <-ch
		//log.Info("Signal received:", sig)
		switch sig {
		case syscall.SIGHUP:
			log.Info("get sighup sighup")
		case syscall.SIGINT:
			log.Info("get SIGINT ,exit!")
			os.Exit(1)
		case syscall.SIGUSR1:
			log.Info("usr1")
			//close the net
			lis.Close()
			log.Info("close connect")
			if _, _, err := syscall.StartProcess(os.Args[0], os.Args, &attr); err != nil {
				check(err)
			}
			//exit current process.
			return
		case syscall.SIGUSR2:
			log.Info("usr2 ")
		}
	}
}
