package main

import (
	"flag"
)

var configPath string

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//传入配置文件
	//var configPath string
	flag.StringVar(&configPath,
		"conf",
		"conf/eggs.conf", "Path of the TOML configuration of the eggs.conf ,default conf/eggs.conf .")
	flag.Parse()

	config, err := ConfigFromFile(configPath)
	if err != nil {
		panic(err.Error())
	}

	//日志配置
	LoadConf(config.Loging.File)

	//启动rest服务
	restservice.LaunchRestServer(configPath)

	//抓取信号
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
