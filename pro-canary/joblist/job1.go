package main

import (
	"fmt"
	"time"
)

type JobId string
type JobData string
type WorkerFactory func() Worker

type WorkerConfig struct {
	Name    JobId
	Factory WorkerFactory
	Count   int // 需要启动的 worker 数量
}

// 所有的任务都会读取 Context 的内容，所以这个结构会很大。
// 当它变得过于复杂的时候需要重构，不过这就不是现在讨论的问题了。
type Context struct {
	Jobs []JobId

	// 各种可能被用到的字段
	Data   JobData
	Fetch  string
	Build  string
	Packer string
}

// 任务调度器
type Dispatcher struct {
	done        chan bool
	jobChannels map[JobId]*JobChannels
}

type JobChannels struct {
	input  chan *Context
	output chan *Context
}

// Worker 的接口
type Worker interface {
	Work(input <-chan *Context, output chan<- *Context)
}

// 各种 worker的struct
type FetchWorker struct{}
type BuildWorker struct{}
type PackerWorker struct{}

func main() {
	fmt.Println("starting...")
	//初始调度器
	dispatcher := NewDispatcher()

	// 这里用来演示通过网络异步收到 job 的情况
	go func() {
		job1 := []JobId{"fetch", "build", "packer"}
		job2 := []JobId{"fetch", "packer"} // 跳过 bar
		job3 := []JobId{"packer", "build"} // 逆序
		job4 := []JobId{"fetch", "packer"} // 跳过 bar
		job5 := []JobId{"fetch", "packer"} // 跳过 bar

		// 执行任务，每个任务可以带一个自定义数据，现在先简单用 string，未来应该根据设计

		dispatcher.Send(job1, "taskid:1")
		dispatcher.Send(job2, "taskid:2")
		dispatcher.Send(job3, "taskid:3")
		dispatcher.Send(job4, "taskid:4")
		dispatcher.Send(job5, "taskid:5")

		time.Sleep(time.Second)

		//如果是http server 不能stop
		dispatcher.Stop()
	}()

	dispatcher.Start()
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		done: make(chan bool),
	}
}

var workerConfig = []*WorkerConfig{
	&WorkerConfig{"fetch", NewFetchWorker, 5},
	&WorkerConfig{"build", NewBuildWorker, 2},
	&WorkerConfig{"packer", NewPackerWorker, 10},
}

func (d *Dispatcher) Start() {
	d.jobChannels = make(map[JobId]*JobChannels)

	// 启动足够数量的 worker
	for _, config := range workerConfig {
		channels := &JobChannels{
			input:  make(chan *Context),
			output: make(chan *Context),
		}
		d.jobChannels[config.Name] = channels

		for i := 0; i < config.Count; i++ {
			worker := config.Factory()
			go worker.Work(channels.input, channels.output)
		}
	}

	// 做输入输出的调度工作
	for _, channels := range d.jobChannels {
		go d.monitor(channels.output)
	}

	<-d.done
}

//输入输出控制
func (d *Dispatcher) monitor(output <-chan *Context) {
	for ctx := range output {
		go d.dispatch(ctx)
	}
}

func (d *Dispatcher) dispatch(ctx *Context) {
	// 所有任务都完成了
	if len(ctx.Jobs) == 0 {
		fmt.Println("all job is done! Name:", ctx.Data, "Data:", *ctx)
		return
	}

	// 把 ctx 放入对应的任务队列，开始执行任务
	job := ctx.Jobs[0]
	ctx.Jobs = ctx.Jobs[1:]
	channels := d.jobChannels[job]
	channels.input <- ctx
}

func (d *Dispatcher) Stop() {
	d.done <- true
}

func (d *Dispatcher) Send(jobs []JobId, data JobData) {
	// 首先初始化一个上下文
	ctx := &Context{
		Jobs: jobs,
		Data: data,
	}

	// 开始派发任务
	d.dispatch(ctx)
}

func NewFetchWorker() Worker {
	return &FetchWorker{}
}

func NewBuildWorker() Worker {
	return &BuildWorker{}
}

func NewPackerWorker() Worker {
	return &PackerWorker{}
}

func (fetch *FetchWorker) Work(input <-chan *Context, output chan<- *Context) {
	for ctx := range input {
		fmt.Println("Worker fetch: current job name is", ctx.Data)
		ctx.Fetch = "fetch-done"
		output <- ctx
	}
}

func (build *BuildWorker) Work(input <-chan *Context, output chan<- *Context) {
	for ctx := range input {
		fmt.Println("Worker build: current job name is", ctx.Data)
		ctx.Build = "build-done"
		output <- ctx
	}
}

func (packer *PackerWorker) Work(input <-chan *Context, output chan<- *Context) {
	for ctx := range input {
		fmt.Println("Worker packer: current job name is", ctx.Data)
		ctx.Packer = "packer-done"
		output <- ctx
	}
}
