package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	//    "io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	//"sync"
	"text/template"
	//"math/rand"
	"database/sql"
	log "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/config"
	"runtime"
	"time"
)

const AppVersion = "Version 2.0.20141121"

var (
	usage = `test`
)

type Job1 struct {
	jobname Info
	results chan<- Result
}

type Job2 struct {
	jobname Info
	results chan<- Result
}

type Job3 struct {
	jobname Info
	results chan<- Result
}


type Result struct {
	JobId    int64
	resultcode int
	resultinfo string
}

type Info struct {
	ID   int64
	Name string
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//获取相关信息
	fetch := NewFetch()

	//并发调度处理
	fetch.doRequest()

	//尾部显示
	pdo.displayEnd()
}

//#################
//主体main的相关函数
//获取和判断相关信息

func N() *Pdo {

	flag.Parse()

	var info *Info
	return info

}

//####################
//并发调度过程
//处理job对列
//并发调度开始
func (pdo *Pdo) doRequest() {
	jobs := make(chan Job, pdo.Concurrent)
	results := make(chan Result, len(pdo.JobList))
	done := make(chan struct{}, pdo.Concurrent)

	go pdo.addJob(jobs, pdo.JobList, results)

	for i := 0; i < pdo.Concurrent; i++ {
		go pdo.doJob(done, jobs)
	}

	go pdo.sysSignalHandle()

	go pdo.awaitCompletion(done, results, pdo.Concurrent)

	pdo.processResults(results)
}

//添加job
func (pdo *Pdo) addJob(jobs chan<- Job, jobnames []HostList, results chan<- Result) {
	for num, jobname := range jobnames {
		jobs <- Job{jobname, results}
		//第一个任务暂停
	}
	close(jobs)
}

//处理job
func (pdo *Pdo) doJob(done chan<- struct{}, jobs <-chan Job) {

	for job := range jobs {
		pdo.Do(&job)
		time.Sleep(pdo.TimeInterval)
	}
	done <- struct{}{}
}

//job完成状态
func (pdo *Pdo) awaitCompletion(done <-chan struct{}, results chan Result, works int) {
	for i := 0; i < works; i++ {
		<-done
	}
	close(results)
}

//job处理结果
func (pdo *Pdo) processResults(results <-chan Result) {
	//0 success
	//1 fail
	//2 time over killed
	//3 time over kill failed

	jobnum := 1
	success := 0
	fail := 0
	overtime := 0

	for result := range results {
		switch result.resultcode {
		case 0:
			CreateAppendFile(pdo.SuccessFile, result.jobname)
			success++
		case 2:
			fmt.Printf("[%d/%d] %s \033[1;31m [Time Over KILLED]\033[0m.\n", jobnum, pdo.JobTotal, result.jobname)
			CreateAppendFile(pdo.FailFile, result.jobname)
			overtime++
		case 3:
			fmt.Printf("[%d/%d] %s \033[1;31m [KILLED FAILED]\033[0m.\n", jobnum, pdo.JobTotal, result.jobname)
			CreateAppendFile(pdo.FailFile, result.jobname)
			overtime++
		default:
			if pdo.OutputWay != "row" {
				fmt.Printf("[%d/%d] %s \033[1;31m [FAILED]\033[0m.\n", jobnum, pdo.JobTotal, result.jobname)
			}
			CreateAppendFile(pdo.FailFile, result.jobname)
			fail++
		}
		fmt.Println(result.resultinfo)
		jobnum++
	}
	fmt.Printf("[INFO] Total: %d ; Success: %d ; Failed: %d ; OverTime: %d\n", pdo.JobTotal, success, fail, overtime)
}

//####################
//公共调用函数
//错误检查
func checkErr(i int, err error) {
	if err != nil {
		switch i {
		case 1:
			log.Critical(err)
		case 2:
			log.Warn(err)
		default:
			log.Info(err)
		}
	}
	log.Flush()
}

//检查文件是否存在.
func checkExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//信号处理
func (pdo *Pdo) sysSignalHandle() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Warnf("Ctrl+c,recode fail list to "+pdo.FailFile+" ,signal:%s", sig)
			for x := pdo.JobFinished - 1; x < pdo.JobTotal; x++ {
				CreateAppendFile(pdo.FailFile, pdo.JobList[x].Host+" "+pdo.JobList[x].Path)
			}
			os.Exit(0)
		}
	}()
}

//具体job处理过程
func (pdo *Pdo) Do(job *Job) {

	pdo.JobFinished++
	var out, outerr bytes.Buffer
	var cmd *exec.Cmd

	//默认是命令行执行
	//本地脚本执行与远程执行
		cmd = exec.Command("ssh", "-q", "-xT", "-o", "PasswordAuthentication=no", "-o", "StrictHostKeyChecking=no", "-o", "ConnectTimeout=3", job.jobname.Host, "cd", job.jobname.Path, "&&", pdo.CmdLastString)
	}
	//命令执行
	stdout, err := cmd.StdoutPipe()
	checkErr(2, err)
	stderr, err := cmd.StderrPipe()
	checkErr(2, err)
	err = cmd.Start()
	checkErr(2, err)

	done := make(chan error)

	go func() {
		done <- cmd.Wait()
	}()

	//线程控制执行时间
	select {
	case <-time.After(pdo.TimeWait):
		//超时被杀时
		if err := cmd.Process.Kill(); err != nil {
			//超时被杀失败
			job.results <- Result{jobstring, 3, "Killed..."}
			checkErr(2, err)
		}
		<-done
		job.results <- Result{jobstring, 2, "Time over ,Killed..."}
		//记录失败job
	case err := <-done:
		if err != nil {
			//完成返回失败
			job.results <- Result{jobstring, 1, outerr.String()}
		} else {
			//完成返回成功
			if pdo.OutputWay != "row" {
				//如果是行显示就隐藏
				job.results <- Result{jobstring, 0, out.String()}
			}
		}
	}

}
