package render

import (
	"errors"
	"geoserver/internal/loader"
	"os"
	"strconv"
	"sync"
	"time"
)

func Tiler(layer loader.LayerGD, z, x, y int) error {
	filePath := "../resource/cache/" + layer.Name + "/" + strconv.Itoa(z) + "/" + strconv.Itoa(x) + "/"
	fileName := strconv.Itoa(y) + "." + layer.TileExt
	file, err := os.Open(filePath + fileName)
	if os.IsNotExist(err) {
		resultChan := make(chan Result, 1)
		MakeTask(layer, z, x, y, resultChan, nil)
		select {
		case result := <-resultChan:
			{
				if result.err != nil {
					return result.err
				}
			}
		case <-time.After(10 * time.Minute):
			return errors.New("Истекло время ожидания воркера")
		}
		file, err = os.Open(filePath + fileName)
		if err != nil {
			return err
		}
	}
	file.Close()
	return nil
}

func MakeTask(layer loader.LayerGD, z, x, y int, resultChan chan Result, wg *sync.WaitGroup) {
	filePath := "../resource/cache/" + layer.Name + "/" + strconv.Itoa(z) + "/" + strconv.Itoa(x) + "/"
	fileName := strconv.Itoa(y) + "." + layer.TileExt
	Tasks <- Task{
		layer:    layer,
		filePath: filePath,
		fileName: fileName,
		x:        x,
		y:        y,
		z:        z,
		result:   resultChan,
		wg:       wg,
	}
}
