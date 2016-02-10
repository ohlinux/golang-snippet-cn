// 基本结构：
//   Task：当前任务共享的上下文，任务通过上下文交换数据，一个任务可分为很多的工作（Work）
//   Dispatcher：任务队列管理器，负责创建 Task 并使用合适的 Worker 来处理数据
//   Worker：任务的抽象接口
//   XXXWorker：各个具体的任务处理逻辑
//   WorkerBench：一个 Worker 池，确保当前正在运行的 Worker 数量不超过限制
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type WorkerId string
type TaskData string
type WorkerFactory func() Worker

type WorkerConfig struct {
	Name    WorkerId
	Factory WorkerFactory
	Count   int // 需要启动的 worker 数量
}

// 所有的任务都会读取 Task 的内容，所以这个结构会很大。
// -1 不需要 0 需要 1 完成 2 失败
type Task struct {
	// 各种可能被用到的字段
	Data   []ModuleInfo
	Tag    TaskData
	TaskId int
	Fetch  string
	Build  string
	Pack   string
}

type ModuleInfo struct {
	AppId      int
	OriginId   int //module id or conf id
	Version    string
	SrcType    string //获取方式scm git svn
	Category   string //类型:conf bin data
	SrcPath    string //来源类型: Source : scm
	Compress   int    //处理方式: true压缩了,需要解压,false 不需要解压.
	DeployPath string //部署位置: DeployPath : /
	SavaPath   string //woker之间传递的目录名.
	FileName   string //保存的文件名
}

// 任务调度器
type Dispatcher struct {
	done          chan bool
	workerBenches map[WorkerId]Worker
}

// 用来创建 Worker，并限制同时工作的 Worker 总数。
type WorkerBench struct {
	throttle chan bool
	factory  WorkerFactory
}

// Worker 的接口
type Worker interface {
	Work(*Task)
}

// 各种 worker
type FetchWorker struct {
}
type BuildWorker struct {
}
type PackWorker struct {
}

func (api *Api) InitDispatcher() {
	dispatcher := NewDispatcher()
	api.DP = dispatcher
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		done:          make(chan bool),
		workerBenches: make(map[WorkerId]Worker, 1),
	}
}

//worker配置
var workerConfig = []*WorkerConfig{
	&WorkerConfig{"fetch", NewFetchWorker, 10},
	&WorkerConfig{"build", NewBuildWorker, 10},
	&WorkerConfig{"pack", NewPackWorker, 5},
}

//调度器启动
func (d *Dispatcher) Start() {
	//workerBenches := make(map[WorkerId]Worker)

	// 初始化 WorkerBench
	for _, config := range workerConfig {
		d.workerBenches[config.Name] = NewWorkerBench(config.Factory, config.Count)
	}

	//d.workerBenches = workerBenches

	<-d.done
}

func (d *Dispatcher) Stop() {
	d.done <- true
}

// 对指定的数据 data 执行一系列工作。
func (d *Dispatcher) Exec(works []ModuleInfo, data TaskData) {
	go d.exec(works, data)
}

func (d *Dispatcher) exec(works []ModuleInfo, data TaskData) {
	// 首先初始化一个上下文

	var bench Worker
	var wg sync.WaitGroup
	//开始分解任务到Module/work 一个任务的fetch和build 结束才能进行pack操作,使用waitgroup.
	//task.Data的数据有写冲突.
	for _, work := range works {
		wg.Add(1)

		go func(w ModuleInfo) {
			//局部变量
			task := &Task{
				Tag:  data,
				Data: []ModuleInfo{w},
			}
			bench = d.workerBenches[WorkerId("fetch")]
			bench.Work(task)

			bench = d.workerBenches[WorkerId("build")]
			bench.Work(task)

			wg.Done()
		}(work)
	}

	wg.Wait()

	task := &Task{
		Tag: data,
	}
	task.Data = works
	bench = d.workerBenches[WorkerId("pack")]
	bench.Work(task)
}

// 初始化一个 WorkerBench，默认标记所有 Worker 都为空闲。
func NewWorkerBench(factory WorkerFactory, count int) *WorkerBench {
	throttle := make(chan bool, count)

	for i := 0; i < count; i++ {
		throttle <- true
	}

	return &WorkerBench{
		throttle: throttle,
		factory:  factory,
	}
}

// 开始执行一项任务。
func (b *WorkerBench) Work(task *Task) {
	<-b.throttle
	worker := b.factory()
	worker.Work(task)
	b.throttle <- true
}

func NewFetchWorker() Worker {
	return &FetchWorker{}
}

func NewBuildWorker() Worker {
	return &BuildWorker{}
}

func NewPackWorker() Worker {
	return &PackWorker{}
}

func (fetch *FetchWorker) Work(task *Task) {
	fmt.Println("Worker fetch:", "Begin:", time.Now(), task.Data, task.Tag)

	task.Fetch = "fetch-done"
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	fmt.Println("Worker fetch: ", "END:", time.Now(), task.Data, task.Tag)

}

func (build *BuildWorker) Work(task *Task) {
	fmt.Println("Worker build: ", "Begin:", time.Now(), task.Data, task.Tag)
	if err := Build(task.Data); err != nil {
		task.Build = "build-failed"
	} else {
		task.Build = "build-done"
	}
	//time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	fmt.Println("Worker build: ", "END:", time.Now(), task.Data, task.Tag)

}

func (pack *PackWorker) Work(task *Task) {
	fmt.Println("Worker pack:  ", "Begin:", time.Now(), task.Data, task.Tag)
	task.Pack = "pack-done"
	time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
	fmt.Println("Worker pack:  ", "END:", time.Now(), task.Data, task.Tag)

}

//模拟request结构
//type Request struct {
//	Id         int64        `json:"id"` //task id
//	ModuleList []PackModule //需要packer的包列表
//}

//模拟main的请求过程与任务分配.
//func main() {
//	fmt.Println("starting...")
//	//初始化调度器
//	dispatcher := NewDispatcher()
//
//	// 这里用来演示通过网络异步收到 work 的情况
//
//	go func() {
//		//进来一批request 之后,按module进行分解
//		r0 := `{"Id":1,"ModulePackModule":[{"Name":"m1","Version":1},{"Name":"m2","Version":2}]}`
//		r1 := `{"Id":2,"ModulePackModule":[{"Name":"m33","Version":1},{"Name":"m22","Version":2},{"Name":"m1","Version":1}]}`
//		r2 := `{"Id":3,"ModulePackModule":[{"Name":"m44","Version":1},{"Name":"m222","Version":2}]}`
//		r3 := `{"Id":4,"ModulePackModule":[{"Name":"m66","Version":1},{"Name":"m2222","Version":2},{"Name":"m1","Version":1}]}`
//		requestPackModule := []string{r0, r1, r2, r3}
//
//		for i, request := range requestPackModule {
//			//处理单个request=task
//			res := &Request{}
//			json.Unmarshal([]byte(request), &res)
//
//			// 执行分配任务，每个任务可以带一个自定义数据，现在先简单用 string
//			fmt.Println(request, TaskData("task"+strconv.Itoa(i)))
//			dispatcher.Exec(res.ModulePackModule, TaskData("workid:"+strconv.Itoa(i)))
//		}
//
//		time.Sleep(100 * time.Second)
//		dispatcher.Stop()
//	}()
//
//	//启动调度器,按不同的channel配置启动不同的channel与worker
//	dispatcher.Start()
//}
