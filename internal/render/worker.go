package render

import (
	"fmt"
	"geoserver/internal/config"
	"geoserver/internal/db/models"
	"geoserver/internal/translator"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

type Task struct {
	layer              models.Layer
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

		// coef := math.Pow(2, float64(task.z))
		// maxSize := min(task.layer.Width, task.layer.Height)
		// xFloat := float64(task.x)
		// yFloat := float64(task.y)
		// maxSizeFloat := float64(maxSize)
		// readSize := maxSizeFloat / coef
		// if xFloat*readSize >= float64(task.layer.Width) || yFloat*readSize >= float64(task.layer.Height) {

		// 	if task.result != nil {
		// 		task.result <- Result{isSuccess: false, err: errors.New("Выход за границы")}
		// 		close(task.result)
		// 	}
		// 	if task.wg != nil {
		// 		task.wg.Done()
		// 	}
		// 	continue
		// }
		// options := []string{"-srcwin",
		// 	fmt.Sprintf("%d", int(xFloat*readSize)),
		// 	fmt.Sprintf("%d", int(yFloat*readSize)),
		// 	fmt.Sprintf("%d", int(readSize)),
		// 	fmt.Sprintf("%d", int(readSize)),
		// 	"-outsize",
		// 	fmt.Sprintf("%d", task.layer.TileSize),
		// 	fmt.Sprintf("%d", task.layer.TileSize)}
		// err = gdal.ConvertTile(
		// 	"../resource/map/"+task.layer.SourcePath,
		// 	task.filePath+task.fileName,
		// 	options,
		// )
		minX, minY, maxX, maxY := translator.WebMercarator(task.x, task.y, task.z)

		// Формируем команду
		args := []string{
			"-s_srs", task.layer.Projection,
			"-t_srs", "EPSG:3857",
			"-te",
			fmt.Sprintf("%f", minX),
			fmt.Sprintf("%f", minY),
			fmt.Sprintf("%f", maxX),
			fmt.Sprintf("%f", maxY),
			"-ts",
			strconv.Itoa(task.layer.TileSize),
			strconv.Itoa(task.layer.TileSize),
			"-r", "lanczos",
			"-of", "PNG",
			"-co", "COMPRESS=DEFLATE",
			"-co", "ZLEVEL=6",
			"-overwrite",
			"../resource/map/" + task.layer.SourcePath,
			task.filePath + task.fileName,
		}

		cmd := exec.Command("gdalwarp", args...)
		err = cmd.Run()
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
