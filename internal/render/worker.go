package render

import (
	"errors"
	"fmt"
	"geoserver/internal/config"
	"geoserver/internal/db/models"
	"math"
	"os"
	"os/exec"
)

type Task struct {
	layer              models.Layer
	z, x, y            int
	filePath, fileName string
	result             chan Result
}

type Result struct {
	isSuccess bool
	err       error
}

var Tasks chan Task

func InitWorkers() {
	Tasks = make(chan Task)
	for i := 0; i < config.Configs.WORKER_COUNT; i++ {
		go renderWorker(Tasks)
	}
}

func renderWorker(tasks <-chan Task) {
	fmt.Println("Start processing")
	for task := range tasks {
		err := os.MkdirAll(task.filePath, 0755)
		if err != nil {
			task.result <- Result{isSuccess: false, err: err}
			close(task.result)
			continue
		}
		_, err = os.Create(task.filePath + task.fileName)
		if err != nil {
			task.result <- Result{isSuccess: false, err: err}
			close(task.result)
			continue
		}

		coef := math.Pow(2, float64(task.z))
		maxSize := min(task.layer.Width, task.layer.Height)
		xFloat := float64(task.x)
		yFloat := float64(task.y)
		maxSizeFloat := float64(maxSize)
		readSize := maxSizeFloat / coef
		if xFloat*readSize >= float64(task.layer.Width) || yFloat*readSize >= float64(task.layer.Height) {
			task.result <- Result{isSuccess: false, err: errors.New("Выход за границы")}
			close(task.result)
			os.Remove(task.filePath + task.fileName)
			continue
		}
		cmd := exec.Command("gdal_translate", "-srcwin",
			fmt.Sprintf("%d", int(xFloat*readSize)),
			fmt.Sprintf("%d", int(yFloat*readSize)),
			fmt.Sprintf("%d", int(readSize)),
			fmt.Sprintf("%d", int(readSize)),
			"-outsize",
			fmt.Sprintf("%d", task.layer.TileSize),
			fmt.Sprintf("%d", task.layer.TileSize),
			"../resource/map/"+task.layer.SourcePath,
			task.filePath+task.fileName)
		err = cmd.Run()
		if err != nil {
			task.result <- Result{isSuccess: false, err: err}
			close(task.result)
			os.Remove(task.filePath + task.fileName)
			continue
		}
		task.result <- Result{isSuccess: true, err: nil}
		close(task.result)
	}
}
