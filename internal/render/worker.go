package render

import (
	"fmt"
	"geoserver/internal/config"
	"geoserver/internal/loader"
	"os"
	"sync"
)

type Task struct {
	layer              loader.LayerGD
	z, x, y            int
	filePath, fileName string
	result             chan Result
	wg                 *sync.WaitGroup
}

type Result struct {
	isSuccess bool
	err       error
}

var Tasks chan Task

func InitWorkers() {
	Tasks = make(chan Task, config.Configs.WORKER_COUNT*10)
	for i := 0; i < config.Configs.WORKER_COUNT; i++ {
		go renderWorker(Tasks)
	}
}

func renderWorker(tasks <-chan Task) {
	fmt.Println("Start processing")
	for task := range tasks {
		err := os.MkdirAll(task.filePath, 0755)
		if err != nil {

			if task.result != nil {
				task.result <- Result{isSuccess: false, err: err}
				close(task.result)
			}
			if task.wg != nil {
				task.wg.Done()
			}
			continue
		}
		err = WarpRender(task)
		if err != nil {

			if task.result != nil {
				task.result <- Result{isSuccess: false, err: err}
				close(task.result)
			}
			os.Remove(task.filePath + task.fileName)
			if task.wg != nil {
				task.wg.Done()
			}
			continue
		}
		if task.result != nil {
			task.result <- Result{isSuccess: true, err: nil}
			close(task.result)
		}
		if task.wg != nil {
			task.wg.Done()
		}

	}
}
