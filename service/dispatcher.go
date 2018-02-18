package service

import (
	"fmt"
	"sync"
	"github.com/labstack/gommon/log"
)

type Dispatcher struct {
	wg     *sync.WaitGroup
	logger *log.Logger

	TaskQueue  chan Task
	WorkerPool chan chan Worker
}

func NewDispatcher() *Dispatcher {
	taskQueue := make(chan Task)
	workerPool := make(chan chan Worker)
	return &Dispatcher{
		wg:         &sync.WaitGroup{},
		TaskQueue:  taskQueue,
		WorkerPool: workerPool,
	}
}

func (l *Dispatcher) Start() {
	defer l.wg.Done()
	fmt.Println("Dispatcher started.")
}

func (l *Dispatcher) Stop() {
	fmt.Println("Dispatcher exited.")
}

func (l *Dispatcher) dispatch() {

}
