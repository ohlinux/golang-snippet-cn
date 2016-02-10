// 一个简单的并发可控、任务可随意拼接的任务队列实现。
// 仅作概念演示用，细节不要纠结。
//
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
	"strconv"
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
// 当它变得过于复杂的时候需要重构，不过这就不是现在讨论的问题了。
type Task struct {
	// 各种可能被用到的字段
	Data   TaskData
	Foo    string
	Bar    string
	Player string
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
type FooWorker struct{}
type BarWorker struct{}
type PlayerWorker struct{}

func main() {
	fmt.Println("starting...")
	dispatcher := NewDispatcher()

	// 这里用来演示通过网络异步收到 work 的情况
	go func() {
		testWorks := [][]WorkerId{
			[]WorkerId{"foo", "bar", "player"},
			[]WorkerId{"foo", "bar", "player"},
			[]WorkerId{"foo", "bar", "player"},
			[]WorkerId{"foo", "bar", "player"},
			[]WorkerId{"foo", "bar", "player"},
			[]WorkerId{"foo", "bar", "player"},
			[]WorkerId{"foo", "player"}, // 跳过 bar
			[]WorkerId{"foo", "player"}, // 跳过 bar
			[]WorkerId{"foo", "player"}, // 跳过 bar
			[]WorkerId{"foo", "player"}, // 跳过 bar
			[]WorkerId{"foo", "player"}, // 跳过 bar
			[]WorkerId{"foo", "player"}, // 跳过 bar
			[]WorkerId{"bar", "foo"},    // 逆序
			[]WorkerId{"bar", "foo"},    // 逆序
			[]WorkerId{"bar", "foo"},    // 逆序
			[]WorkerId{"bar", "foo"},    // 逆序
			[]WorkerId{"bar", "foo"},    // 逆序
		}

		// 执行任务，每个任务可以带一个自定义数据，现在先简单用 string
		for i, works := range testWorks {
			//fmt.Println(works, TaskData("work"+strconv.Itoa(i)))
			dispatcher.Exec(works, TaskData("work"+strconv.Itoa(i)))
		}

		time.Sleep(time.Second)
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
	&WorkerConfig{"foo", NewFooWorker, 2},
	&WorkerConfig{"bar", NewBarWorker, 1},
	&WorkerConfig{"player", NewPlayerWorker, 3},
}

func (d *Dispatcher) Start() {
	workerBenches := make(map[WorkerId]Worker)

	// 初始化 WorkerBench
	for _, config := range workerConfig {
		workerBenches[config.Name] = NewWorkerBench(config.Factory, config.Count)
	}

	d.workerBenches = workerBenches

	<-d.done
}

func (d *Dispatcher) Stop() {
	d.done <- true
}

// 对指定的数据 data 执行一系列工作。并行task.
func (d *Dispatcher) Exec(works []WorkerId, data TaskData) {
	go d.exec(works, data)
}

func (d *Dispatcher) exec(works []WorkerId, data TaskData) {
	// 首先初始化一个上下文
	task := &Task{
		Data: data,
	}

	// 开始执行所有的任务
	fmt.Println("开始执行:", works)
	for _, work := range works {
		bench := d.workerBenches[work]
		bench.Work(task)
	}
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

func NewFooWorker() Worker {
	return &FooWorker{}
}

func NewBarWorker() Worker {
	return &BarWorker{}
}

func NewPlayerWorker() Worker {
	return &PlayerWorker{}
}

func (foo *FooWorker) Work(task *Task) {
	fmt.Println("Worker foo: current work name is", task.Data)
	task.Foo = "foo-done"
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
}

func (bar *BarWorker) Work(task *Task) {
	fmt.Println("Worker bar: current work name is", task.Data)
	task.Bar = "bar-done"
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
}

func (player *PlayerWorker) Work(task *Task) {
	fmt.Println("Worker player: current work name is", task.Data)
	task.Player = "player-done"
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
}
