package runner

import (
	"time"
	"errors"
	"os"
	"os/signal"
)

//管理程序的生命周期

//在runner 指定的时间完成任务，并且在收到中断信号的时候停止任务

type Runner struct {
	//通道  报告  从系统发出的信号
	interrupt chan os.Signal
	//通信 报告  任务已经完成
	complete chan error
	//timeout 报告任务超时
	timeout <-chan time.Time

	//tasks 持有一组以索引为顺序依次执行的的函数

	tasks []func(int)
}

var ErrTimeOut = errors.New("received timeout")
var ErrInterrupt = errors.New("received interrupt")

//准备一个新的Runner
func New(d time.Duration) *Runner{
	return &Runner{
		interrupt:make(chan os.Signal,1),
		complete:make(chan error),
		timeout:time.After(d),
	}
}

//Add一个Runner ,这个任务接受一个int作为参数的函数

func (r *Runner) Add(tasks ...func(int)){
	r.tasks = append(r.tasks,tasks...)
}

func(r *Runner) start() error {
	//希望接受所有中断信号
	signal.Notify(r.interrupt,os.Interrupt)
	//用不同的go routine 在执行不同的任务
	go func() {
		r.complete <- r.run()
	}()
	select {
	case err:= <-r.complete:
		return err
	case := <- r.timeout:
		return ErrInterrupt
	}
}

//Run执行每一个注册的任务

func (r *Runner) run() error {
	for id, task := range r.tasks {
		if r.getInterrupt() {
			return ErrInterrupt
		}
		task(id)
	}
	return nil
}

func (r *Runner) getInterrupt() bool{
	select {
	case <-r.interrupt:
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}
